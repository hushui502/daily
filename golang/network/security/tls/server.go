package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	cert, err := tls.LoadX509KeyPair("hu.name.pem", "private.pem")
	if err != nil {
		return
	}
	config := tls.Config{Certificates:[]tls.Certificate{cert}}

	now := time.Now()
	config.Time = func() time.Time {
		return now
	}

	config.Rand = rand.Reader

	service := "0.0.0.0:8080"

	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		return
	}
	fmt.Println("Listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Accepted")
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		fmt.Println("Trying to read")
		n, err := conn.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		_, err = conn.Write(buf[0:n])
		if err != nil {
			return
		}
	}
}