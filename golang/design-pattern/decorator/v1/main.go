package main

import "fmt"

func decorator(f func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("started")
		f(s)
		fmt.Println("done")
	}
}

func hello(s string) {
	fmt.Println(s)
}

func main() {
	decorator(hello)("hello~~")
}
