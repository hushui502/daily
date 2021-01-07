package main

import (
	"errors"
	"net"
	"net/rpc"
)

// A program will call a function with a list of parameters, and on completion of the function call will
// have a set of return values.These values may be the function value,of if address have been passed as
// parameters then the comments of those address might have been changed.

// The remote procedure call is an attempt to bring this style of programming into the network world.
// Thus a client will make what looks to it like a normal procedure call.The client-side will package this
// into a network message and transfer it to the server.The server will unpack this and turn it back into
// a procedure call on the server side.The result of this call will be packaged up for return to the client

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B

	return nil
}

func main() {
	arith := new(Arith)
	rpc.Register(arith)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	if err != nil {
		// ...
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		// ...
	}

	// this works
	// rpc.Accept(listener)

	// and so does this
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}