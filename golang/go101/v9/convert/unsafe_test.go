package main

import (
	"log"
	"reflect"
	"unsafe"
)

func stringtoslicebyte(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func slicebytetostring(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: bh.Data,
		Len:  bh.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

func main() {
	s1 := "abc"
	b1 := []byte("def")
	copy(b1, s1)
	log.Println(s1, b1)

	s := "hello"
	b2 := stringtoslicebyte(s)
	log.Println(b2)
	// b2[0] = byte(99) unexpected fault address

	b3 := []byte("test")
	s3 := slicebytetostring(b3)
	log.Println(s3)
}
