package main

import (
	"fmt"
	"net"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// close connection on exit
	defer conn.Close()

	var buf [512]byte
	for {
		// read upto 512 bytes
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println(string(buf[0:]))
		// write the n bytes read
		_, err = conn.Write(buf[0:n])
		if err != nil {
			return
		}
	}
}
