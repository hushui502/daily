package main

import (
	"fmt"
	"net/rpc"
	"os"
)

//type Args struct {
//	A, B int
//}
//
//type Quotient struct {
//	Quo, Rem int
//}

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	serverAddress := os.Args[1]

	client, err := rpc.DialHTTP("tcp", serverAddress+":8080")
	if err != nil {
		fmt.Println("dialing error: ", err)
	}

	args := Args{18, 4}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		fmt.Println("Arith error: ", err)
	}
	fmt.Printf("Arith: %d * %d = %d\n", args.A, args.B, reply)

	var quot Quotient
	err = client.Call("Arith.Divide", args, quot)
	if err != nil {
		// ,,,
	}
	fmt.Printf("Arith: %d / %d = %d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)


}
