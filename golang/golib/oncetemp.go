package golib

import (
	"sync"
	"sync/atomic"
)

var (
	instance *singleton
	once     Once
)

func NewInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
