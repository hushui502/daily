package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(inCh <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range inCh {
			out <- n * n
			time.Sleep(time.Second)
		}
	}()

	return out
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(cs))
	collect := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}

	for _, c := range cs {
		go collect(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func use() {
	in := producer(1, 2, 4)

	c1 := square(in)
	c2 := square(in)

	for res := range merge(c1, c2) {
		fmt.Println(res)
	}
}
