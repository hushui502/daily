package ants

import (
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	capacity int32
	running int32
	workers workerArray
	state int32
	lock sync.Locker
	cond *sync.Cond
	workerCache sync.Pool
	blockingNum int
	options *Options
}

func (p *Pool) periodicallyPurge() {
	heartbeat := time.NewTicker(p.options.ExpiryDuration)
	defer heartbeat.Stop()

	for range heartbeat.C {
		if atomic.LoadInt32(&p.state) == CLOSED {
			break
		}

		p.lock.Lock()
		expiredWorkers := p.workers.retrieveExpiry(p.options.ExpiryDuration)
		p.lock.Unlock()

		for i := range expiredWorkers {
			expiredWorkers[i].task <- nil
		}

		if p.Running() == 0 {
			p.cond.Broadcast()
		}
	}
}

func NewPool(size int, options ...Option) (*Pool, error) {
	opts := loadOptions(options...)

	if size <= 0 {
		size = -1
	}

	if expiry := opts.ExpiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		opts.ExpiryDuration = DefaultCleanIntervalTime
	}

	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}

	p := &Pool{
		capacity:	int32(size),
		lock:	NewSpinLock(),
		options: opts,
	}
	p.workerCache.New = func() interface{} {
		return &goWorker{
			pool:p,
			task:make(chan func(), workerChanCap),
		}
	}
	if p.options.PreAlloc {
		if size == -1 {
			return nil, ErrInvalidPreAllocSize
		}
		p.workers = newWorkerArray(loopQueueType, size)
	} else {
		p.workers = newWorkerArray(stackType, 0)
	}

	p.cond = sync.NewCond(p.lock)

	go p.periodicallyPurge()

	return p, nil
}

func (p *Pool) Submit(task func()) error {
	if atomic.LoadInt32(&p.state) == CLOSED {
		return ErrPoolClosed
	}
	var w *goWorker
	if w = p.retrieveWorker(); w == nil {
		return ErrPoolOverload
	}
	w.task <- task
	return nil
}

func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

func (p *Pool) Free() int {
	return p.Cap() - p.Running()
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

func (p *Pool) Tune(size int) {
	if capacity := p.Cap(); capacity == -1 || size <= 0 || size == capacity || p.options.PreAlloc {
		return
	}

	atomic.StoreInt32(&p.capacity, int32(size))
}

func (p *Pool) Release() {
	atomic.StoreInt32(&p.state, CLOSED)
	p.lock.Lock()
	p.workers.reset()
	p.lock.Unlock()
}

func (p *Pool) Reboot() {
	if atomic.CompareAndSwapInt32(&p.state, CLOSED, OPENED) {
		go p.periodicallyPurge()
	}
}

func (p *Pool) incRunning() {
	atomic.AddInt32(&p.running, 1)
}

func (p *Pool) decRunning() {
	atomic.AddInt32(&p.running, -1)
}

func (p *Pool) retrieveWorker() (w *goWorker) {
	spawnWorker := func() {
		w = p.workerCache.Get().(*goWorker)
		w.run()
	}
	p.lock.Lock()

	w = p.workers.detach()
	if w != nil {
		p.lock.Unlock()
	} else if capacity := p.Cap(); capacity == -1 {
		p.lock.Unlock()
		spawnWorker()
	} else {
		if p.options.NonBlocking {
			p.lock.Unlock()
			return
		}
	Reentry:
		if p.options.MaxBlockingTasks != 0 && p.blockingNum >= p.options.MaxBlockingTasks {
			p.lock.Unlock()
			return
		}
		p.blockingNum++
		p.cond.Wait()
		p.blockingNum--
		if p.Running() == 0 {
			p.lock.Unlock()
			spawnWorker()
			return
		}
		w = p.workers.detach()
		if w == nil {
			goto Reentry
		}

		p.lock.Unlock()
	}
	return
}

func (p *Pool) revertWorker(worker *goWorker) bool {
	if capacity := p.Cap(); (capacity > 0 && p.Running() > capacity) || atomic.LoadInt32(&p.state) == CLOSED {
		return false
	}
	worker.recycleTime = time.Now()
	p.lock.Lock()

	err := p.workers.insert(worker)
	if err != nil {
		p.lock.Unlock()
		return false
	}

	p.cond.Signal()
	p.lock.Unlock()
	return true
}