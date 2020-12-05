package main

import (
	"fmt"
	"time"
)

func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, v := range nums {
			out <- v
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

func use() {
	in := producer(1, 2, 3)
	ch := square(in)
	for res := range ch {
		fmt.Println(res)
	}
}
