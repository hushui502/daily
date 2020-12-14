package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	a1 := []int{0, 1, 2, 3}

	for i, v := range a1 {
		fmt.Println(i, v)

		if i == 0 {
			a1 = a1[2:]
		}
	}
	fmt.Println(a1)
}

func printSlice(s []int) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Printf("header=%+v len=%d cap=%d %v\n", sh, len(s), cap(s), s)
}
