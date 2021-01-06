package main

import (
	"bytes"
	"io"
	"net"
	"time"
)

//func main() {
//	if len(os.Args) != 2 {
//		os.Exit(1)
//	}
//	service := os.Args[1]
//	udpAddr, err := net.ResolveUDPAddr("udp4", service)
//	if err != nil {
//		return
//	}
//	conn, err := net.DialUDP("udp", nil, udpAddr)
//  // net.Dial("tcp", service)
//	if err != nil {
//		return
//	}
//
//	var buf [512]byte
//	n, err := conn.Read(buf[0:])
//	if err != nil {
//		return
//	}
//	fmt.Println(string(buf[0:n]))
//
//	os.Exit(1)
//}

func main() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}

	for {
		handleClient(conn)
	}
}

func listenC(network, service string) {
	// network "tcp"
	// service ":1200"
	listener, err := net.Listen(network, service)
	if err != nil {
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleListenClient(conn)
	}
}

func handleListenClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		_, err = conn.Write(buf[0:n])
		if err != nil {
			return
		}
	}
}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return result.Bytes(), nil
}
