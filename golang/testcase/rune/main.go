package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var a = "a蛤"
	// [97 232 155 164]
	fmt.Println([]byte(a))

	for _, r := range a {
		// 输出2次4，因为这里值得是rune，每个rune其实相当于一个int32，4个字节
		fmt.Println(unsafe.Sizeof(r))
	}

	// a占1个字节，蛤占3个字节 ===> 1+3 = 4
	fmt.Println(len(a))

}

// 0x4F60在0x0800-0xFFFF之间，UTF-8使用3字节模板：1110xxxx 10xxxxxx 10xxxxxx

