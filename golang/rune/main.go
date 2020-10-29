package main

import (
	"fmt"
	"time"
)

func main() {
	skipIfStillRunning()
}

func skipIfStillRunning() {
	m := make(chan struct{}, 1)
	m <- struct{}{}

	for {
		select {
		case v := <-m:
			fmt.Println("running")
			time.Sleep(time.Second * 2)
			m <- v
		default:
			fmt.Println("skip")
		}
	}
}