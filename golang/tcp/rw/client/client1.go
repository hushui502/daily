package main

import (
	"log"
	"net"
	"os"
	"time"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
	}
	defer conn.Close()
	log.Println("dial ok")

	data := os.Args[1]
	conn.Write([]byte(data))

	time.Sleep(time.Second * 100)
}
