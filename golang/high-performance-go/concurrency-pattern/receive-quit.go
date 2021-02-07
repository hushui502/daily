package main

import "fmt"

func main() {
	quit := make(chan string)
	ch := generator("hello", quit)
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}

	<-quit
	fmt.Printf("Generator says: %s", <-quit)
}

func generator(msg string, quit chan string) <-chan string{
	ch := make(chan string)

	go func() {
		for {
			select {
			case ch <- fmt.Sprintf("%s", msg):
			case <-quit:
				quit <- "see you"
				return
			}
		}
	}()

	return ch
}
