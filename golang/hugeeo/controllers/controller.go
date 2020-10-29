package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type ControllerInfo struct {
}

func main() {
	//a1 := make([]byte, 1024)
	file, _ := os.Open("hugeeo/unsafe_test.go")

	buf := bufio.NewReader(file)
	for {
		c, b, err := buf.ReadLine()
		println(b)
		if err != nil {
			if err == io.EOF {
				println("end-----")
				break
			}
		}
		fmt.Println(string(c))

	}
}
