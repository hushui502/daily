package main

import (
	"io"
	"log"
	"net"
	"os"
	"github.com/integrii/flaggy"
	"time"
)

/**
 * Very fast and simple go program that transparently proxies any raw TCP connection from one TCP socket to another.
 * The Wormhole destination is the same every time and specified by flags at startup time.
 */

// Listener is the listen address for incoming connections
var Listener string

// Destination is the target that connections will be forward to
var Destination string

func init() {
	Listener = os.Getenv("LISTENER")
	Destination = os.Getenv("DESTINATION")

	flaggy.String(&Listener, "f", "from", "The address and port that wormhole should listen on.  Connections enter here.")
	flaggy.String(&Destination, "t", "to", "Specifies the address and port that wormhole should redirect TCP connections to.  Connections exit here.")
	flaggy.Parse()
}

func main() {
	log.Println("Starting up. Forwarding connections from ", Listener, "to ", Destination)

	l, err := net.Listen("tcp", Listener)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// wait for a connection.
		c, err := l.Accept()
		if err != nil {
			log.Println("Error establshing incoming connection: ", err)
			continue
		}
		log.Println("Client connected from ", c.RemoteAddr())

		// handle the connection in a goroutine
		go wormhole(c)
	}
}

// wormhole opens a wormhole from the client connection
// to user the specified destination
func wormhole(c net.Conn) {
	defer c.Close()
	log.Println("Opening a wormhole from ", c.RemoteAddr())

	start := time.Now()

	// connect to the destination tcp port
	destConn, err := net.Dial("tcp", Destination)
	if err != nil {
		log.Println("Error connecting to destination port: ", err)
		return
	}
	defer destConn.Close()
	log.Println("Wormhole open from ", c.RemoteAddr())

	// low speed
	go func() { io.Copy(c, destConn) }()

	io.Copy(destConn, c)

	end := time.Now()
	// use lib instead of (+ - * /)
	duration := end.Sub(start)
	log.Println("Closing wormhole from", c.RemoteAddr(), "after", duration)
}