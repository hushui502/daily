package main

import (
	"fmt"
	"sync"
	"time"
)

type Reporter struct {
	worker int
	messages chan string
	wg sync.WaitGroup
	closed bool
}

func NewReporter(worker, buffer int) *Reporter {
	return &Reporter{
		worker:   worker,
		messages: make(chan string, buffer),
	}
}

func (r *Reporter) run(stop <-chan struct{}) {
	// 在mian over执行，避免往一个关闭channel里发信息
	go func() {
		<-stop
		r.shutdown()
	}()

	for i := 0; i < r.worker; i++ {
		r.wg.Add(1)
		go func() {
			for msg := range r.messages {
				// 模拟消费message
				time.Sleep(2 * time.Second)
				fmt.Printf("report: %s\n", msg)
			}
			r.wg.Done()
		}()
	}

	r.wg.Wait()
}

func (r *Reporter) shutdown() {
	r.closed = true
	close(r.messages)
}

func (r *Reporter) report(data string) {
	if r.closed {
		return
	}
	r.messages <- data
}