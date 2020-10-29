package main

import "fmt"

func main() {
	v := []int{1, 2, 3}
	for _, d := range v {
		v = append(v, d)
	}
	fmt.Println(v)
}
