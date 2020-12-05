package main

import "sync"

func main() {
	useLock()
}

type Lock struct {
	c chan struct{}
}

func NewLock() Lock {
	var l Lock
	l.c = make(chan struct{}, 1)
	l.c <- struct{}{}
	return l
}

func (l Lock) Lock() bool {
	lockResult := false
	select {
	case <-l.c:
		lockResult = true
	default:
	}
	return lockResult
}

func (l Lock) UnLock() {
	l.c <- struct{}{}
}

var counter int

func useLock() {
	var l = NewLock()
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !l.Lock() {
				println("lock failed")
				return
			}
		}()
		counter++
		println("counter ", counter)
		l.UnLock()
	}

	wg.Wait()
}