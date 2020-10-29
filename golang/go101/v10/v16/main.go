package main

import (
	"fmt"
	"github.com/andybalholm/cascadia"
	"io"
	"log"
	"net"
	"time"
)

func forward(src net.Conn, network, address string, timeout time.Duration) {
	defer src.Close()
	dst, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		log.Println("dial err")
		return
	}
	defer dst.Close()

	cpErr := make(chan error)

	go cp(cpErr, src, dst)
	go cp(cpErr, dst, src)

	select {
	case err = <-cpErr:
		if err != nil {
			log.Println("copy err: %v", err)
		}
	}
	log.Printf("disconnect: %s", src.RemoteAddr())

}

func cp(c chan error, w io.Writer, r io.Reader) {
	_, err := io.Copy(w, r)
	c <- err
	fmt.Println("cp end")
}
