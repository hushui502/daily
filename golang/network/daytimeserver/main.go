package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println(err)
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

		daytime := time.Now().String()
		conn.Write([]byte(daytime))
		conn.Close()
	}
}
