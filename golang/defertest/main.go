package main

import "runtime"

var i int

func main() {
	runtime.GOMAXPROCS(1)

	go func() {
		panic("already call")
	}()

	for {
		i++
	}
}
