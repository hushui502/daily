package ringbuffer

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken

	// https://github.com/golang/go/issues/13086
	// https://medium.com/a-journey-with-go/go-mutex-and-starvation-3f4f4e75ad50
	mutexStarving
	mutexWaiterShift = iota
)

// Mutex is a Locker
type Mutex struct {
	sync.Mutex
}

func (m *Mutex) TryLock() bool {
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}

	new := old | mutexLocked

	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, new)
}

func (m *Mutex) Count() int {
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	v = v >> mutexWaiterShift
	v = v + (v & mutexLocked)

	return int(v)
}

func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))

	return state&mutexLocked == mutexLocked
}

func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))

	return state&mutexWoken == mutexWoken
}

func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))

	return state&mutexStarving == mutexStarving
}
