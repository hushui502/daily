package fishpool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Errors
var (
	ErrPoolNotRunning = errors.New("the pool is not running")
	ErrJobNotFunc     = errors.New("generic worker not given a func()")
	ErrWorkerClosed   = errors.New("worker was closed")
	ErrJobTimedOut    = errors.New("job request timed out")
)

type Worker interface {
	// Process will synchronously perform a job and return the result.
	Process(interface{}) interface{}

	// BlockUntilReady is called before each job is processed and must block
	// the calling goroutine until the worker is ready to process the next job.
	BlockUntilReady()

	// Interrupt is called when a job is canceled.
	Interrupt()

	// Terminate is called when a worker is removed from the processing pool
	// and is responsible for cleaning up any held resources.
	Terminate()
}

// closureWorker is a minimal Worker implementation that simply wraps a
// func(interface{}) interface{}
type closureWorker struct {
	processor func(interface{}) interface{}
}

func (w *closureWorker) Process(payload interface{}) interface{} {
	return w.processor(payload)
}

func (w *closureWorker) BlockUntilReady() {}
func (w *closureWorker) Interrupt()       {}
func (w *closureWorker) Terminate()       {}

// callbackWorker is a minimal Worker implementation that attempts to cast
// each job into func() and either calls it if successful or returns
// ErrJobNotFunc.
type callbackWorker struct{}

func (w *callbackWorker) Process(payload interface{}) interface{} {
	f, ok := payload.(func())
	if !ok {
		return ErrJobNotFunc
	}
	f()
	return nil
}

func (w *callbackWorker) BlockUntilReady() {}
func (w *callbackWorker) Interrupt()       {}
func (w *callbackWorker) Terminate()       {}

// Pool is a struct that manages a collection of workers, each with their own
// goroutine. The Pool can initialize, expand, compress and close the workers,
// as well as processing jobs with the workers synchronously.
type Pool struct {
	queuedJobs int64

	ctor    func() Worker
	workers []*workerWrapper
	reqChan chan workRequest

	workerMut sync.Mutex
}

// New creates a new Pool for workers that starts with n workers.
// You must provide a constructor function that new Worker types and
// when you change the size of the pool the constructor will be called
// to create each new Worker.
func New(n int, ctor func() Worker) *Pool {
	p := &Pool{
		ctor:    ctor,
		reqChan: make(chan workRequest),
	}
	p.SetSize(n)

	return p
}

// NewFunc creates a new Pool of workers where each worker will
// process using the provided func
func NewFunc(n int, f func(interface{}) interface{}) *Pool {
	return New(n, func() Worker {
		return &closureWorker{
			processor: f,
		}
	})
}

// NewFunc creates a new Pool of workers where each worker will process
// using the provided func.
func NewCallback(n int) *Pool {
	return New(n, func() Worker {
		return &closureWorker{}
	})
}

// Process will use the Pool to process a payload and synchronously return
// resutl. Process can be called safely by any goroutines, but will panic if
// the Pool has been stoped.
func (p *Pool) Process(payload interface{}) interface{} {
	atomic.AddInt64(&p.queuedJobs, 1)

	request, open := <-p.reqChan
	if !open {
		panic(ErrPoolNotRunning)
	}

	request.jobChan <- payload

	payload, open = <-request.retChan
	if !open {
		panic(ErrWorkerClosed)
	}

	atomic.AddInt64(&p.queuedJobs, -1)

	return payload
}

// ProcessTimed will use the Pool to process a payload and synchronously return
// the result. If the timeout occurs before the job has finished the worker will
// be interrupted and ErrJobTimedOut will be returned. ProcessTimed can be
// called safely by any goroutines.
func (p *Pool) ProcessTimed(payload interface{}, timeout time.Duration) (interface{}, error) {
	atomic.AddInt64(&p.queuedJobs, 1)
	defer atomic.AddInt64(&p.queuedJobs, -1)

	tout := time.NewTimer(timeout)

	var request workRequest
	var open bool

	select {
	case request, open = <-p.reqChan:
		if !open {
			return nil, ErrPoolNotRunning
		}
	case <-tout.C:
		return nil, ErrJobTimedOut
	}

	select {
	case request.jobChan <- payload:
	case <-tout.C:
		request.interruptFunc()
		return nil, ErrJobTimedOut
	}

	select {
	case payload, open = <-request.retChan:
		if !open {
			return nil, ErrWorkerClosed
		}
	case <-tout.C:
		request.interruptFunc()
		return nil, ErrJobTimedOut
	}

	tout.Stop()
	return payload, nil
}

// ProcessCtx will use the Pool to process a payload and synchronously return
// the result. If the context cancels before the job has finished the worker will
// be interrupted and ErrJobTimedOut will be returned. ProcessCtx can be
// called safely by any goroutines.
func (p *Pool) ProcessCtx(ctx context.Context, payload interface{}) (interface{}, error) {
	atomic.AddInt64(&p.queuedJobs, 1)
	defer atomic.AddInt64(&p.queuedJobs, -1)

	var request workRequest
	var open bool

	select {
	case request, open = <-p.reqChan:
		if !open {
			return nil, ErrPoolNotRunning
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case request.jobChan <- payload:
	case <-ctx.Done():
		request.interruptFunc()
		return nil, ctx.Err()
	}

	select {
	case payload, open = <-request.retChan:
		if !open {
			return nil, ErrWorkerClosed
		}
	case <-ctx.Done():
		request.interruptFunc()
		return nil, ctx.Err()
	}

	return payload, nil
}

func (p *Pool) QueueLength() int64 {
	return atomic.LoadInt64(&p.queuedJobs)
}

// SetSize changes the total number of wokers in the Pool.
// This can be called by any goroutine at any time unless the Pool
// has been stopped, in which case a panic will occur.
func (p *Pool) SetSize(n int) {
	p.workerMut.Lock()
	defer p.workerMut.Unlock()

	lworkers := len(p.workers)
	if lworkers == n {
		return
	}

	// add extra workers if N > len(workers)
	for i := lworkers; i < n; i++ {
		p.workers = append(p.workers, newWorkerWrapper(p.reqChan, p.ctor()))
	}

	// stop all workers > N
	for i := n; i < lworkers; i++ {
		p.workers[i].stop()
	}

	// wait for all workers > N to stop
	for i := n; i < lworkers; i++ {
		p.workers[i].join()
	}

	// remove stopped workers from slice
	p.workers = p.workers[:n]
}

// GetSize returns the current size of the pool
func (p *Pool) GetSize() int {
	p.workerMut.Lock()
	defer p.workerMut.Unlock()

	return len(p.workers)
}

// Close will terminate all workers and close the job channel of this pool
func (p *Pool) Close() {
	p.SetSize(0)
	close(p.reqChan)
}
