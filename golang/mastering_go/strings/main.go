package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	r := strings.NewReader("hello")
	fmt.Println("r length ", r.Len())
	b := make([]byte, 1)
	for {
		n, err := r.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		fmt.Printf("read %s bytes: %d\n", b, n)
	}

	fmt.Println("-------")
	s := strings.NewReader("This is an error!\n")
	n, err := s.WriteTo(os.Stderr)
	if err != nil {
		return
	}
	fmt.Printf("wrote %d bytes to os.stderr\n", n)

}
