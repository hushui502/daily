package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type RecursiveMutex struct {
	sync.Mutex
	owner int64		// 当前持有锁的goroutine id
	recursion int32	// 当前持有锁的goroutine i
}

func (m *RecursiveMutex) Lock() {
	gid := GoID()
	// 如果当前持有锁的goroutine就是这次调⽤的goroutine,说明是重⼊
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的goroutine第⼀次调⽤，记录下它的goroutine id,调⽤次数加1
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := GoID()
	// ⾮持有锁的goroutine尝试释放锁，错误的使⽤
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	// 调⽤次数减1
	m.recursion--
	if m.recursion != 0 { // 如果这个goroutine还没有完全释放，则直接返回
		return
	}
	// 此goroutine最后⼀次调⽤，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

func GoID() int64 {
	var buf [64]byte

	n := runtime.Stack(buf[:], false)

	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine"))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}

	return int64(id)
}

func main() {
	//mu := RecursiveMutex{}
	//
	//mu.Lock()
	//mu.Lock()
	//mu.Lock()
	//println(mu.recursion)

	c := sync.NewCond(&sync.Mutex{})
	var ready int

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)

			c.L.Lock()
			ready++
			c.L.Unlock()

			log.Println("ok=============")

			c.Broadcast()
		}(i)
	}
	sync.Once{}
	c.L.Lock()
	for ready != 10 {
		c.Wait()
		log.Println("wake up once")
	}
	c.L.Unlock()

	println("over!!!!!!!!!")
}



