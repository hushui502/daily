package main

import (
	"fmt"
	"net"
)

// the type TCPAddr is a structure containing an IP and a port
// type TCPAddr struct {
// 		IP IP
//		Port int
// }

func main() {
	// network is "tcp" or "udp"
	// address is "www.baidu.com:89" {hostname:port}
	tcpaddr, err := net.ResolveTCPAddr("tcp", "www.baidu.com:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*tcpaddr)
}

// go run main.go
// {180.101.49.11 80 }