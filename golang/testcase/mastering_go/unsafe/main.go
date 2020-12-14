package main


import (
	"fmt"
	"unsafe"
)

func main() {
	var value int64 = 6
	var p1 = &value
	var p2 = (*int32)(unsafe.Pointer(p1))
	fmt.Println(*p2)

	fmt.Println("==========")
	array := [...]int{1, 2, 3, 4}
	pointer := &array[0]
	fmt.Println(*pointer)
	memoryAddress := uintptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(array[0])
	for i := 0; i < len(array)-1; i++ {
		pointer = (*int)(unsafe.Pointer(memoryAddress))
		fmt.Println(*pointer)
		memoryAddress = uintptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(array[0])
	}
}
