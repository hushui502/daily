package main

import (
	"fmt"
	"net"
	"os"
)

const (
	DIR = "DIR"
	CD = "CD"
	PWD = "PWD"
)

func main() {
	service := "0.0.0.0:1202"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		s := string(buf[0:n])
		if s[0:2] == CD {
			chdir(conn, s[3:])
		} else if s[0:3] == DIR {
			dirList(conn)
		} else if s[0:3] == PWD {
			pwd(conn)
		}
	}
}

func chdir(conn net.Conn, s string) {
	if os.Chdir(s) == nil {
		conn.Write([]byte("OK"))
	} else {
		conn.Write([]byte("ERROR"))
	}
}

func pwd(conn net.Conn) {
	s, err := os.Getwd()
	if err != nil {
		conn.Write([]byte(""))
		return
	}
	
	conn.Write([]byte(s))
}

func dirList(conn net.Conn) {
	defer conn.Write([]byte("\r\n"))

	dir, err := os.Open(".")
	if err != nil {
		return
	}

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}

	for _, name := range names {
		conn.Write([]byte(name + "\r\n"))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("err is ", err)
		os.Exit(1)
	}
}

// ---------------pseudocode-------------------
//state = login
//while true
//	read line
//	swich state
//		case login:
//			get name from line
//			get password from line
//			if name and password verified
//				write successed
//				state = file_transfer
//			else
//				write failed
//				state = login
//		case file_transfer
//				if line.startwith cd
//					get dir from line
//					if changedir dir ok
//						write successed
//						state = file_transfer
//					else
//						write failed
//						state = file_transfer
//						...


// -----------------summary----------------------
// Building any application requires design decisions before you start writing code.
// For distributed applications you should have a wider range of decisions to make compared to standalone systems.