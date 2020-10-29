package main

import (
	"fmt"
	"unsafe"
)

type A struct {
	arr   [2]int8  // 2 数组主要取决于item
	slice []int8   // 24 切片看SliceHeader
	bl    bool     // 1
	sl    []int16  // 24
	ptr   *int64   // 8 指针类型64位是8个字节，32位是4个字节
	st    struct { // 16 struct取决于字段
		str string // string其实就是一个指针和一个int表示的len，所以8+8=16
	}
	m map[string]int16 // map其实就是一个指针引用
	i interface{}      // interface实际是一个类型指针一个数据指针，9+8=16
}

func main() {
	a := A{}
	// 2 1
	fmt.Println(unsafe.Sizeof(a.arr), unsafe.Alignof(a.arr))
	// 24 8
	fmt.Println(unsafe.Sizeof(a.arr), unsafe.Alignof(a.arr))
	// 1 1
	fmt.Println(unsafe.Sizeof(a.bl), unsafe.Alignof(a.bl))
	// 24 8
	fmt.Println(unsafe.Sizeof(a.sl), unsafe.Alignof(a.sl))
	// 8 8
	fmt.Println(unsafe.Sizeof(a.ptr), unsafe.Alignof(a.ptr))
	// 16 8
	fmt.Println(unsafe.Sizeof(a.st), unsafe.Alignof(a.st))
	// 8 8
	fmt.Println(unsafe.Sizeof(a.m), unsafe.Alignof(a.m))
	// 16 8
	fmt.Println(unsafe.Sizeof(a.i), unsafe.Alignof(a.i))
}
