package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	b := []byte("hello")
	s := string(b)
	t := string(b)
	// 824634236680
	fmt.Println(pointer(s))
	// 824634236648
	fmt.Println(pointer(t))

}

func pointer(s string) uintptr {
	p := unsafe.Pointer(&s)
	h := *(*reflect.StringHeader)(p)
	return h.Data
}

func intern(m map[string]string, b []byte) string {
	c, ok := m[string(b)]
	if ok {
		return c
	}

	s := string(b)
	m[s] = s
	return s
}
