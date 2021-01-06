package main

import (
	"crypto/tls"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	service := os.Args[1]

	conn, err := tls.Dial("tcp", service, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for n := 0; n < 10; n++ {
		conn.Write([]byte("hello " + string(n+48)))

		var buf [512]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(string(buf[0:n]))
	}
	os.Exit(1)
}

// Conclusion
// Security is a huge area in itself, and in this chapter we have barely touched on it.However,the major concepts
// have been covered.What has not been stressed is how much security needs to be built into the design phase.
// Security as an afterthought is nearly always a failure.
