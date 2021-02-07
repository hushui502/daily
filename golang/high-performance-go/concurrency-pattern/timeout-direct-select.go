package main

import (
	"fmt"
	"time"
)

func main() {
	ch := generator("hello")
	timeout := time.After(4 * time.Second)
	for {
		select {
		case s := <-ch:
			fmt.Println(s)
		case <-timeout:
			return
		}
	}
}

func generator(msg string) <-chan string { // returns receive-only channel
	ch := make(chan string)
	go func() { // anonymous goroutine
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Second)
		}
	}()
	return ch
}