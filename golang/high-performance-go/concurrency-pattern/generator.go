package main

import (
	"fmt"
	"time"
)

func main() {
	ch := generator("hello")
	// 对channel进行消费
	for i := 0; i < 3; i++ {
		fmt.Println(<-ch)
	}
}

// returns receive-only channel
// 这里其实是生成了一个channel对象，这个channel对象只能<-ch这样使用
func generator(msg string) <-chan string {
	ch := make(chan string)
	// 模拟一个goroutine调用
	go func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Second)
		}
	}()
	return ch
}
