package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println(
		unsafe.Sizeof(interface{}(0)),
		unsafe.Sizeof(struct{}{}),
	)
}
