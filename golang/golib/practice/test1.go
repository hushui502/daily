package main

import (
	"sync"
	"time"
)

func main() {
	done := make(chan int, 10)

	for i := 0; i < 10; i++ {
		go func() {
			println("hello")
			done <- 1
		}()
	}

	for i := 0; i < cap(done); i++ {
		<-done
	}

	wait()

	ch := make(chan int, 3)
	go Producer(3, ch)
	go Producer(3, ch)
	go Consumer(ch)

	time.Sleep(5 * time.Second)
}

func wait() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			println("world")
			wg.Done()
		}()
	}

	wg.Wait()
}

func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- i * factor
	}
}

func Consumer(in <-chan int) {
	for v := range in {
		println(v)
	}
}
