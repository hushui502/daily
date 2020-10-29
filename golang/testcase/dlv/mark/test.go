package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var byteArray = [10]byte{'a1', 'b', 'c'}
	var p *int64 = (*int64)(unsafe.Pointer(&byteArray[1]))
	fmt.Println(*p)
}
