package main

import (
	"fmt"
	"time"
)

/*
Hello 0
World 0
Hello 1
World 1
Hello 2
World 2
Hello 3
World 3
*/
func main() {
	ch1 := generator("Hello")
	ch2 := generator("World")
	for i := 0; i < 4; i++ {
		// 每一个迭代的i都是一样的，因为会阻塞
		fmt.Println(<-ch1)
		fmt.Println(<-ch2)
	}
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
