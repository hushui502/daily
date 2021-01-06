package main

import (
	"fmt"
	"net"
	"unicode/utf16"
)

// Managing character sets and encodings

func main() {
	str := "百度地图"
	runes := utf16.Encode([]rune(str))
	fmt.Println(runes)
	println(len([]rune(str))) // 4
	println(len(str)) // 12

	//for _, v := range str {
	//	fmt.Printf("%x\n", v)
	//}
}

const BOM = '\ufffe'

func writeShorts(conn net.Conn, shorts []uint16) {
	var bytes [2]byte

	bytes[1] = BOM & 255
	_, err := conn.Write(bytes[0:])
	if err != nil {
		return
	}
	for _, v := range shorts {
		bytes[0] = byte(v >> 8)
		bytes[1] = byte(v & 255)
		_, err = conn.Write(bytes[0:])
		if err != nil {
			return
		}
	}
}

func readShorts(conn net.Conn) []uint16 {
	var buf [512]byte

	n, err := conn.Read(buf[0:2])
	if err != nil {
		return nil
	}
	for {
		m, err := conn.Read(buf[n:])
		if m == 0 || err != nil {
			break
		}
		n += m
	}

	var shorts []uint16
	shorts = make([]uint16, n / 2)
	if buf[0] == 0xff && buf[1] == 0xfe {
		// big endian
		for i := 2; i < n; i += 2 {
			shorts[i/2] = uint16(buf[i]) << 8 + uint16(buf[i+1])
		}
	} else if buf[1] == 0xff && buf[0] == 0xfe {
		// little endian
		for i := 2; i < n; i += 2 {
			shorts[i/2] = uint16(buf[i+1]) << 8 + uint16(buf[i])
		}
	} else {
		// unknown byte error
		fmt.Println("Unknown order")
	}

	return shorts
}

// Conclusion
// There has not been much code in this chapter.Instead, there have been some of the concepts of a very complex area.
// It is up to you:if you want to assume everyone speaks US English then the world is simple.
// But if you want to your applications to be usable by the rest of the world, the you need to pay attention to these complexities.