package main

import "fmt"

func reverse1(str []byte) {
	i := 0
	k := len(str) - 1

	for i < k {
		str[i], str[k] = str[k], str[i]
		i++
		k--
	}
}

func reverse2(str []byte) {
	i := 0
	k := 0

	reverse1 := func(str []byte, begin, end int) {
		for begin < end {
			str[begin], str[end] = str[end], str[begin]
			begin++
			end--
		}
	}

	//reverse1(str, 0, len(str)-1)
	for i = 0; i < len(str); i++ {
		if str[i] == ' ' {
			reverse1(str, k, i - 1)
			k = i + 1
		}
	}

	fmt.Println(string(str))
}

func main() {
	str := "hu fan"
	reverse2([]byte(str))

	
}
