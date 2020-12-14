package main

import (
	"github.com/EDDYCJY/redis-protocol-example/protocol"
	"log"
	"net"
	"os"
)

const (
	Address = "127.0.0.1:6379"
	Network = "tcp"
)

func conn(network, address string) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		log.Fatal("os.args <= 0")
	}

	reqCommand := protocol.GetRequest(args)

	redisConn, err := conn(Network, Address)
	if err != nil {
		log.Fatalf("Conn err : %v", err)
	}
	defer redisConn.Close()

	_, err = redisConn.Write(reqCommand)
	if err != nil {
		log.Fatalf("Conn Write err : %v", err)
	}

	command := make([]byte, 1024)
	n, err := redisConn.Read(command)
	if err != nil {
		log.Fatalf("Conn Read err : %v", err)
	}

	reply, err := protocol.GetReply(command[:n])
	if err != nil {
		log.Fatalf("protocal.Getreply err %v", err)
	}

	log.Printf("Reply: %v", reply)
	log.Printf("Command: %v", string(command[:n]))
}
