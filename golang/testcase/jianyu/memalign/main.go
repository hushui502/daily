package main

import (
	"fmt"
	"sync"
	"time"
)

type Part1 struct {
	a bool
	b int32
	c int8
	d int64
	e byte
}

type Part2 struct {
	c int8
	a bool
	e byte
	b int32
	d int64
}

var wg sync.WaitGroup

func main() {
	userCount := 10
	ch := make(chan int, 5)
	for i := 0; i < userCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range ch {
				fmt.Println(d)
				time.Sleep(time.Second * time.Duration(d))
			}
		}()
	}

	for i := 0; i < 10; i++ {
		ch <- 1
		ch <- 2
	}

	close(ch)
	wg.Wait()
}
