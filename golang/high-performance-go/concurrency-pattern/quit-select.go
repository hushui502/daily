package main

import "fmt"

func main() {
	quit := make(chan struct{})
	ch := generator("hello", quit)
	for i := 0; i < 4; i++ {
		fmt.Println(<-ch)
	}
	quit <- struct{}{}
}

func generator(msg string, quit chan struct{}) <-chan string {
	ch := make(chan string)

	go func() {
		for {
			select {
			case ch <- fmt.Sprintf("%s", msg):
			case <-quit:
				return
			}
		}
	}()

	return ch
}
