package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: netfwd localIp:localPort remoteIp:remotePort")
	}

	localAddr := os.Args[1]
	remoteAddr := os.Args[2]
	local, err := net.Listen("tcp", localAddr)
	if local == nil {
		log.Fatal("cannot listen: %v", err)
	}
	for {
		conn, err := local.Accept()
		if err, ok := err.(net.Error); ok && err.Temporary() {
			continue
		}
		if conn == nil {

		}
		go forward(conn, remoteAddr)
	}
}

func forward(local net.Conn, remoteAddr string) {
	remote, err := net.Dial("tcp", remoteAddr)
	if remote == nil {
		fmt.Fprintf(os.Stderr, "remote dial failed: %v\n", err)
		return
	}
	go io.Copy(local, remote)
	go io.Copy(remote, local)
}
