package main

import "io"

func main() {
	req := make([]byte, 2)
	str := "hufanisagoodman"
	for {
		n, err := io.ReadFull(str, req)
	}
}
