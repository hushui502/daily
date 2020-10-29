package ants

import (
	"runtime"
	"time"
)

type goWorker struct {
	// pool who owns this worker
	pool *Pool

	// task is a job should be done
	task chan func()

	recycleTime time.Time
}

func (w *goWorker) run() {
	w.pool.incRunning()

	go func() {
		defer func() {
			w.pool.decRunning()
			w.pool.workerCache.Put(w)
			if p := recover(); p != nil {
				if ph := w.pool.options.PanicHandler; ph != nil {
					ph(p)
				} else {
					w.pool.options.Logger.Printf("worker exits from a panic %v\n", p)
					var buf [4069]byte
					n := runtime.Stack(buf[:], false)
					w.pool.options.Logger.Printf("worker exits from panic: %s\n", string(buf[:n]))
				}
			}
		}()

		for f := range w.task {
			if f == nil {
				return
			}
			f()
			if ok := w.pool.revertWorker(w); !ok {
				return
			}
		}
	}()
}
