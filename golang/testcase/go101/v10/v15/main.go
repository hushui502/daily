package main

import "fmt"

func main() {
	a := [3]int{1, 2, 3}
	b := (&a)[:]
	println(&a)
	println(b)
	fmt.Println(b)
}
