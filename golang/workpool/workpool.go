// github.com/gammazero/workerpool
package workpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gammazero/deque"
)

const (
	// if works idle for at least this period of time, then stop a worker.
	idleTimeout = 2 * time.Second
)

// WorkerPool is a collection of goroutines, where the number of concurrent
// goroutines processing requests does not exceed the specified maximum.
type WorkerPool struct {
	maxWorkers   int
	taskQueue    chan func()
	workerQueue  chan func()
	stoppedChan  chan struct{}
	stopSignal   chan struct{}
	waitingQueue deque.Deque
	stopLock     sync.Mutex
	stopOnce     sync.Once
	stopped      bool
	waiting      int32
	wait         bool
}

func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := &WorkerPool{
		maxWorkers:  maxWorkers,
		taskQueue:   make(chan func(), 1),
		workerQueue: make(chan func()),
		stopSignal:  make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}

	// start the task dispatcher.
	go pool.dispatch()

	return pool
}

// returns the maximum number of concurrent workers.
func (p *WorkerPool) Size() int {
	return p.maxWorkers
}

// stop the worker pool and wait for only currently running tasks to
// complete. pending tasks that are not currently running are abandoned.
// tasks must not be submitted to the worker pool after calling stop.
func (p *WorkerPool) Stop() {
	p.stop(false)
}

func (p *WorkerPool) StopWait() {
	p.stop(true)
}

func (p *WorkerPool) Stopped() bool {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()

	return p.stopped
}

func (p *WorkerPool) Submit(task func()) {
	if task != nil {
		p.taskQueue <- task
	}
}

func (p *WorkerPool) SubmitWait(task func()) {
	if task == nil {
		return
	}
	doneChan := make(chan struct{})
	p.taskQueue <- func() {
		task()
		close(doneChan)
	}

	<-doneChan
}

func (p *WorkerPool) WaitingQueueSize() int {
	return int(atomic.LoadInt32(&p.waiting))
}

func (p *WorkerPool) Pause(ctx context.Context) {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()

	if p.stopped {
		return
	}
	ready := &sync.WaitGroup{}
	ready.Add(p.maxWorkers)
	for i := 0; i < p.maxWorkers; i++ {
		p.Submit(func() {
			ready.Done()
			select {
			case <-ctx.Done():
			case <-p.stopSignal:
			}
		})
	}

	// wait fro workers to all be paused
	ready.Wait()
}

func (p *WorkerPool) dispatch() {
	defer close(p.stoppedChan)
	timeout := time.NewTimer(idleTimeout)
	var workerCount int
	var idle bool
	var wg sync.WaitGroup

Loop:
	for {
		if p.waitingQueue.Len() != 0 {
			if !p.processWaitingQueue() {
				break Loop
			}
			continue
		}

		select {
		case task, ok := <-p.taskQueue:
			if !ok {
				break Loop
			}
			// get the task to do.
			select {
			case p.workerQueue <- task:
			default:
				// create a new worker, if not at max.
				if workerCount < p.maxWorkers {
					wg.Add(1)
					go startWorker(task, p.workerQueue, &wg)
					workerCount++
				} else {
					// enqueue task to be executed by next available worker.
					p.waitingQueue.PushBack(task)
					atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len()))
				}
			}
			idle = false
		case <-timeout.C:
			if idle && workerCount > 0 {
				if p.killIdleWorker() {
					workerCount--
				}
			}
			idle = true
			timeout.Reset(idleTimeout)
		}
	}

	// if instructed to wait, then run tasks that are already queued.
	if p.wait {
		p.runQueuedTasks()
	}

	// stop all remaining workers as they become ready.
	for workerCount > 0 {
		p.workerQueue <- nil
		workerCount--
	}
	wg.Wait()

	timeout.Stop()
}

// startWorker runs initial task, then starts a worker waiting for more.
func startWorker(task func(), workQueue chan func(), wg *sync.WaitGroup) {
	task()
	go worker(workQueue, wg)
}

// worker executes tasks and stops when it receives a nil task.
func worker(workerQueue chan func(), wg *sync.WaitGroup) {
	for task := range workerQueue {
		if task == nil {
			wg.Done()
			return
		}
		task()
	}
}

// stop tells the dispatcher to exit, and whether or not to complete queued tasks.
func (p *WorkerPool) stop(wait bool) {
	p.stopOnce.Do(func() {
		close(p.stopSignal)
		p.stopLock.Lock()
		p.stopped = true
		p.stopLock.Unlock()
		p.wait = wait
		close(p.taskQueue)
	})

	// who give it?
	<-p.stoppedChan
}

func (p *WorkerPool) processWaitingQueue() bool {
	select {
	case task, ok := <-p.taskQueue:
		if !ok {
			return false
		}
		p.waitingQueue.PushBack(task)
	case p.workerQueue <- p.waitingQueue.Front().(func()):
		// a worker was ready, so gave task to worker
		p.waitingQueue.PopFront()
	}
	atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len()))

	return true
}

func (p *WorkerPool) killIdleWorker() bool {
	select {
	case p.workerQueue <- nil:
		// sent kill signal to worker.
		return true
	default:
		// no ready workers. all if any workers are busy.
		return false
	}
}

func (p *WorkerPool) runQueuedTasks() {
	for p.waitingQueue.Len() != 0 {
		// a worker is ready, so give task to worker.
		p.workerQueue <- p.waitingQueue.PopFront().(func())
		atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len()))
	}
}
