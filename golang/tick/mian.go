package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.Tick(time.Millisecond * 500)

	for range timer {
		fmt.Println("hello")
	}
}
