package main

import (
	"fmt"
	"time"
)

func main() {
	go regular_print("hello")
	fmt.Println("Second print statement!")
	time.Sleep(3 * time.Second)
	// when main returns the goroutines also will be end
	// 一般来讲防止main goroutine结束导致的其他goroutine也结束
	// 我们通常使用一些for select或者信号量，中断信号等来处理
	// 不推荐使用time sleep
	fmt.Println("Third print statement!")
}

func regular_print(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Second)
	}
}


