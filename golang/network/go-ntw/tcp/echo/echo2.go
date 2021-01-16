package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var addr string
	flag.StringVar(&addr, "e", "/tmp/echo2.sock", "service address endpoint")
	flag.Parse()
	text := flag.Arg(0)

	raddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.DialUnix("unix", nil, raddr)
	if err != nil {
		fmt.Println("failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// send text to server
	_, err = conn.Write([]byte(text))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// read response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("failed reading response:", err)
		os.Exit(1)
	}
	fmt.Println(string(buf[:n]))
}
