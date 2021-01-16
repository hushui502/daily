package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var addr string
	flag.StringVar(&addr, "e", ":4040", "service address endpoint")
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		os.Exit(1)
	}

	// announce service using listenTCP which a TCPListener
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("listening at (ctp)", laddr.String())

	for {
		// server accepts tcp conn from client
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("failed to accept conn:", err)
			conn.Close()
			continue
		}
		fmt.Println("connected to: ", conn.RemoteAddr())

		// server handle service that process this conn
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	// read message from conn which client-server
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// write message and send it to client
	w, err := conn.Write(buf[:n])
	if err != nil {
		fmt.Println("failed to write to client:", err)
		return
	}

	// because this is a echo service, so w == n
	if w != n {
		fmt.Println("warning: not all data sent to client")
		return
	}
}