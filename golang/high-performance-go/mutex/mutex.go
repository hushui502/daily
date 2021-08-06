//// Copyright 2009 The Go Authors. All rights reserved.
//// Use of this source code is governed by a BSD-style
//// license that can be found in the LICENSE file.
//
//// Package sync provides basic synchronization primitives such as mutual
//// exclusion locks. Other than the Once and WaitGroup types, most are intended
//// for use by low-level library routines. Higher-level synchronization is
//// better done via channels and communication.
////
//// Values containing the types defined in this package should not be copied.
//package sync
//
//import (
//	"internal/race"
//	"sync/atomic"
//	"unsafe"
//)
//
//func throw(string) // provided by runtime
//
//// A Mutex is a mutual exclusion lock.
//// The zero value for a Mutex is an unlocked mutex.
////
//// A Mutex must not be copied after first use.
//type Mutex struct {
//	state int32
//	sema  uint32
//}
//
//// A Locker represents an object that can be locked and unlocked.
//type Locker interface {
//	Lock()
//	Unlock()
//}
//
//const (
//	mutexLocked = 1 << iota // mutex is locked		持有锁标志
//	mutexWoken										// 唤醒标志
//	mutexStarving									// 饥饿标志
//	mutexWaiterShift = iota
//
//	// Mutex fairness.
//	//
//	// Mutex can be in 2 modes of operations: normal and starvation.
//	// In normal mode waiters are queued in FIFO order, but a woken up waiter
//	// does not own the mutex and competes with new arriving goroutines over
//	// the ownership. New arriving goroutines have an advantage -- they are
//	// already running on CPU and there can be lots of them, so a woken up
//	// waiter has good chances of losing. In such case it is queued at front
//	// of the wait queue. If a waiter fails to acquire the mutex for more than 1ms,
//	// it switches mutex to the starvation mode.
//	//
//	// In starvation mode ownership of the mutex is directly handed off from
//	// the unlocking goroutine to the waiter at the front of the queue.
//	// New arriving goroutines don't try to acquire the mutex even if it appears
//	// to be unlocked, and don't try to spin. Instead they queue themselves at
//	// the tail of the wait queue.
//	//
//	// If a waiter receives ownership of the mutex and sees that either
//	// (1) it is the last waiter in the queue, or (2) it waited for less than 1 ms,
//	// it switches mutex back to normal operation mode.
//	//
//	// Normal mode has considerably better performance as a goroutine can acquire
//	// a mutex several times in a row even if there are blocked waiters.
//	// Starvation mode is important to prevent pathological cases of tail latency.
//	starvationThresholdNs = 1e6
//)
//
//// Lock locks m.
//// If the lock is already in use, the calling goroutine
//// blocks until the mutex is available.
//func (m *Mutex) Lock() {
//	// Fast path: grab unlocked mutex.
//	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
//		if race.Enabled {
//			race.Acquire(unsafe.Pointer(m))
//		}
//		return
//	}
//	// Slow path (outlined so that the fast path can be inlined)
//	// 缓慢之路，尝试⾃旋竞争或饥饿状态下饥饿goroutine竞争
//	m.lockSlow()
//}
//
//func (m *Mutex) lockSlow() {
//	var waitStartTime int64
//	starving := false		// goroutine 饥饿标志
//	awoke := false			// 唤醒标志
//	iter := 0				// 自旋标志
//	old := m.state			// 当前锁的状态
//	for {
//		// Don't spin in starvation mode, ownership is handed off to waiters
//		// so we won't be able to acquire the mutex anyway.
//		// 锁是非饥饿状态，锁还没释放，尝试自旋
//		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
//			// Active spinning makes sense.
//			// Try to set mutexWoken flag to inform Unlock
//			// to not wake other blocked goroutines.
//			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
//				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
//				awoke = true
//			}
//			runtime_doSpin()
//			iter++
//			old = m.state			// 再次获得锁的状态，检查锁是否被释放
//			continue
//		}
//		new := old
//		// Don't try to acquire starving mutex, new arriving goroutines must queue.
//		if old&mutexStarving == 0 {
//			new |= mutexLocked		// 非饥饿状态，加锁
//		}
//		if old&(mutexLocked|mutexStarving) != 0 {
//			new += 1 << mutexWaiterShift	// waiter + 1
//		}
//		// The current goroutine switches mutex to starvation mode.
//		// But if the mutex is currently unlocked, don't do the switch.
//		// Unlock expects that starving mutex has waiters, which will not
//		// be true in this case.
//		if starving && old&mutexLocked != 0 {
//			new |= mutexStarving		// 设置饥饿状态
//		}
//		if awoke {
//			// The goroutine has been woken from sleep,
//			// so we need to reset the flag in either case.
//			if new&mutexWoken == 0 {
//				throw("sync: inconsistent mutex state")
//			}
//			new &^= mutexWoken		 // 新状态清除唤醒标记
//		}
//		// 成功设置新状态
//		if atomic.CompareAndSwapInt32(&m.state, old, new) {
//			// 原来锁的状态已释放，并且不是饥饿状态，正常请求到了锁，返回
//			if old&(mutexLocked|mutexStarving) == 0 {
//				break // locked the mutex with CAS
//			}
//
//			// 处理饥饿状态
//
//			// 如果以前就在队列⾥⾯，加⼊到队列头
//			// If we were already waiting before, queue at the front of the queue.
//			queueLifo := waitStartTime != 0
//			if waitStartTime == 0 {
//				waitStartTime = runtime_nanotime()
//			}
//			// 阻塞等待
//			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
//			// 唤醒之后检查锁是否应该处于饥饿状态
//			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
//			old = m.state
//			// 如果锁已经处于饥饿状态，直接抢到锁，返回
//			if old&mutexStarving != 0 {
//				// If this goroutine was woken and mutex is in starvation mode,
//				// ownership was handed off to us but mutex is in somewhat
//				// inconsistent state: mutexLocked is not set and we are still
//				// accounted as waiter. Fix that.
//				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
//					throw("sync: inconsistent mutex state")
//				}
//				// 加锁并且将waiter数减1
//				delta := int32(mutexLocked - 1<<mutexWaiterShift)
//				if !starving || old>>mutexWaiterShift == 1 {
//					// Exit starvation mode.
//					// Critical to do it here and consider wait time.
//					// Starvation mode is so inefficient, that two goroutines
//					// can go lock-step infinitely once they switch mutex
//					// to starvation mode.
//					delta -= mutexStarving			// 最后⼀个waiter或者已经不饥饿了，清除饥饿标记
//				}
//				atomic.AddInt32(&m.state, delta)
//				break
//			}
//			awoke = true
//			iter = 0
//		} else {
//			old = m.state
//		}
//	}
//
//	if race.Enabled {
//		race.Acquire(unsafe.Pointer(m))
//	}
//}
//
//// Unlock unlocks m.
//// It is a run-time error if m is not locked on entry to Unlock.
////
//// A locked Mutex is not associated with a particular goroutine.
//// It is allowed for one goroutine to lock a Mutex and then
//// arrange for another goroutine to unlock it.
//func (m *Mutex) Unlock() {
//	if race.Enabled {
//		_ = m.state
//		race.Release(unsafe.Pointer(m))
//	}
//
//	// Fast path: drop lock bit.
//	new := atomic.AddInt32(&m.state, -mutexLocked)
//	if new != 0 {
//		// Outlined slow path to allow inlining the fast path.
//		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
//		m.unlockSlow(new)
//	}
//}
//
//func (m *Mutex) unlockSlow(new int32) {
//	if (new+mutexLocked)&mutexLocked == 0 {
//		throw("sync: unlock of unlocked mutex")
//	}
//	if new&mutexStarving == 0 {
//		old := new
//		for {
//			// If there are no waiters or a goroutine has already
//			// been woken or grabbed the lock, no need to wake anyone.
//			// In starvation mode ownership is directly handed off from unlocking
//			// goroutine to the next waiter. We are not part of this chain,
//			// since we did not observe mutexStarving when we unlocked the mutex above.
//			// So get off the way.
//			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
//				return
//			}
//			// Grab the right to wake someone.
//			new = (old - 1<<mutexWaiterShift) | mutexWoken
//			if atomic.CompareAndSwapInt32(&m.state, old, new) {
//				runtime_Semrelease(&m.sema, false, 1)
//				return
//			}
//			old = m.state
//		}
//	} else {
//		// Starving mode: handoff mutex ownership to the next waiter, and yield
//		// our time slice so that the next waiter can start to run immediately.
//		// Note: mutexLocked is not set, the waiter will set it after wakeup.
//		// But mutex is still considered locked if mutexStarving is set,
//		// so new coming goroutines won't acquire it.
//		runtime_Semrelease(&m.sema, true, 1)
//	}
//}
