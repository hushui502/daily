package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	const sample = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"
	// sample len:  8
	fmt.Println("sample len: ", len(sample))

	// bd b2 3d bc 20 e2 8c 98
	for i := 0; i < len(sample); i++ {
		fmt.Printf("%x ", sample[i])
	}

	// bdb23dbc20e28c98
	// bdb23dbc20e28c98
	// "\xbd\xb2=\xbc ⌘"
	// "\xbd\xb2=\xbc \u2318"
	fmt.Printf("%x\n", sample)

	fmt.Printf("%x \n", sample)

	fmt.Printf("%q\n", sample)

	fmt.Printf("%+q\n", sample)

	fmt.Println("======================")

	const placeOfInterest = `⌘`
	// placeOfInterest len:  3
	fmt.Println("placeOfInterest len: ", len(placeOfInterest))

	// plain string:
	// ⌘
	fmt.Println("plain string:")
	fmt.Printf("%s", placeOfInterest)
	fmt.Printf("\n")

	// quoted string:"\u2318"
	fmt.Printf("quoted string:")
	fmt.Printf("%+q", placeOfInterest)
	fmt.Printf("\n")

	// hex bytes:e28c98
	fmt.Printf("hex bytes:")
	for i := 0; i < len(placeOfInterest); i++ {
		fmt.Printf("%x", placeOfInterest[i])
	}
	fmt.Printf("\n")

	fmt.Println("=================")

	const nihongo = "日本語"
	// nihongo len:  9
	fmt.Println("nihongo len: ", len(nihongo))

	// U+65E5 '日' starts at byte position 0
	// U+672C '本' starts at byte position 3
	// U+8A9E '語' starts at byte position 6
	for index, runeValue := range nihongo {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
	}

	// U+65E5 '日' starts at byte position 0
	// U+672C '本' starts at byte position 3
	// U+8A9E '語' starts at byte position 6
	for i, w := 0, 0; i < len(nihongo); i += w {
		runeValue, width := utf8.DecodeRuneInString(nihongo[i:])
		fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
		w = width
	}
}
