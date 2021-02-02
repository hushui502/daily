package main

import (
	"log"
	"sync"
	"time"
)

// sync.Cond条件变量用来协调想要访问共享变量的那些goroutine
// 当共享变量发生变化的时候，它可以用来通知被互斥锁阻塞的goroutine

var done bool

func read(name string, c *sync.Cond) {
	c.L.Lock()
	println("got the lock")
	for !done {
		c.Wait()
	}
	log.Println(name, "starts reading!")
	c.L.Unlock()
}

func write(name string, c *sync.Cond) {
	log.Println(name, "starts writing")
	time.Sleep(2 * time.Second)
	c.L.Lock()
	done = true
	c.L.Unlock()
	log.Println(name, "wakes all")
	// 这里的唤醒是通知并唤醒所有正在wait的cond
	c.Broadcast()
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})

	go read("reader1", cond)
	go read("reader2", cond)
	go read("reader3", cond)

	write("writer", cond)

	time.Sleep(10 * time.Second)
}
