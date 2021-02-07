package main

import "fmt"

func f(left, right chan int) {
	left <- 1 + <-right
}

func main() {
	const n = 10000

	// construct an array of n+1 int channels
	var channels [n + 1]chan int
	for i := range channels {
		channels[i] = make(chan int)
	}

	// wire n goroutines in a chain
	for i := 0; i < n; i++ {
		go f(channels[i], channels[i+1])
	}

	// insert a value into right-hand end
	// 相当于给最后的right ch 传一个值，这样整个f 协程才能跑起来
	go func(c chan<- int) { c <- 1 }(channels[n])

	// get value from the left-hand end
	fmt.Println(<-channels[0])

}
