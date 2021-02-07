package main

import (
	"fmt"
	"time"
)

func main() {
	ch := generator("Hello")
	for i := 0; i < 10; i++ {
		select {
		case s := <-ch:
			fmt.Println(s)
		case <-time.After(time.Second):
			fmt.Println("waited too long")
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