package main

import (
	"encoding/json"
	"fmt"
)

type T struct {
	a int
	b int
	c string
}

type MyString string

func (m MyString) String() string {
	return fmt.Sprintf("mystring = %s\n", string(m))
}

var _ ByteSize

func main() {
	defer func() {
		if err := recover(); err != nil {

		}
	}()
	panic("sds")
}

var _ json.Marshaler = (*json.RawMessage)(nil)

func Min(a ...int) int {
	min := int(^uint(0) >> 1)
	println(min)
	return 1
}

type ByteSize int

func (b ByteSize) String() string {
	switch {
	case b == 1:
		return fmt.Sprintf("1---")
	case b == 2:
		return fmt.Sprintf("2---")
	}
	return fmt.Sprintf("3----")
}
