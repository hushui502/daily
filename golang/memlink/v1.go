package main

import "fmt"

func main() {
	s := []byte("ss")
	str := "hello"
	s = append(s, str...)
	fmt.Println(len(s))
}
