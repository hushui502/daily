package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	var buffer bytes.Buffer
	buffer.Write([]byte("this is "))
	fmt.Fprintf(&buffer, " a good book")
	buffer.WriteTo(os.Stdout)
	buffer.WriteTo(os.Stdout)

	buffer.Reset()
	buffer.Write([]byte("mastering go!"))
	r := bytes.NewReader([]byte(buffer.String()))
	fmt.Println(buffer.String())
	for {
		b := make([]byte, 3)
		n, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		fmt.Printf("read %s bytes: %d\n", b, n)
	}
}
