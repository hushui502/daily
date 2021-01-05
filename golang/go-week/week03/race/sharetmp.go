package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		// for avoiding data race
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
}


