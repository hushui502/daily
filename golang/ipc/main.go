package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
)

func connect() net.Conn {
	conn, err := net.Dial("tcp", "irc.freenode.net:6667")
	if err != nil {
		panic(err)
	}
	return conn
}

func disconnect(conn net.Conn) {
	conn.Close()
}

func logon(conn net.Conn) {
	sendData(conn, "USER TheManWithTheIceCreamVan 8 * :Someone")
	sendData(conn, "NICK TheManWithTheIceCreamVan")
}

func main() {
	conn := connect()
	logon(conn)

	tp := textproto.NewReader(bufio.NewReader(conn))

	for {
		status, err := tp.ReadLine()
		if err != nil {
			panic(err)
		}
		fmt.Println(status)
	}

	disconnect(conn)
}

func sendData(conn net.Conn, message string) {
	fmt.Fprintf(conn, "%s\r\n", message)
}