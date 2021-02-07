package main

import (
	"fmt"
	"time"
)

func main() {
	ch := fanIn(generator("HELLO"), generator("BYE"))

	// 此时就无法保证两个ch的同步性了
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}

func fanIn(ch1, ch2 <-chan string) <-chan string {
	new_ch := make(chan string)
	go func() {
		for {
			new_ch <- <-ch1
		}
	}()
	go func() {
		for {
			new_ch <- <-ch2
		}
	}()

	return new_ch
}

func generator(msg string) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Second)
		}
	}()

	return ch
}
