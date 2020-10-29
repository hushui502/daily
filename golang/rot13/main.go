package main

import "fmt"

func main() {
	message := "uv vagreangvbany fcnpr fgngvba"

	for i := 0; i < len(message); i++ {    // 迭代字符串中的每一个 ASCII 字符
		c := message[i]
		if c >= 'a' && c <= 'z' {    // 只解密英文字母，至于空格和标点符号则保持不变
			c = c + 13
			if c > 'z' {
				c = c - 26
			}
		}
		fmt.Printf("%c", c)
	}
}
