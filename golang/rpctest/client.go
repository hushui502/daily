package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		log.Fatal("dialing: ", err)
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply string
	err = client.Call("HetlloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
}
