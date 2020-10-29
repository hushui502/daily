package main

import (
	"fmt"
	"time"
)

func worker(start chan bool) {
	heartbeat := time.Tick(3 * time.Second)
	for {
		select {
		// do some stuff
		case <-heartbeat:
			fmt.Println("heartbeat")
		}
	}
}

func main() {
	ch := make(chan bool)
	worker(ch)
}
