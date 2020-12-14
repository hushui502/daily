package main

import (
	"fmt"
	"unicode"
)

func main() {
	const r1 = '€'
	fmt.Println("(int32) r1:", r1)
	fmt.Printf("(HEX) r1: %x\n", r1)
	fmt.Printf("(as a String) r1: %s\n", r1)
	fmt.Printf("(as a character) r1: %c\n", r1)

	// is rune?
	const sL= "\x99\x00ab\x50\x00\x23\x50\x29\x9c"
	fmt.Println(len(sL))
	for i := 0; i < len(sL); i++ {
		if unicode.IsPrint(rune(sL[i])) {
			fmt.Printf("%c\n", sL[i])
		} else {
			fmt.Println("not printable")
		}
	}


}
