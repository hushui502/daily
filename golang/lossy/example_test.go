package lossy

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

const packetCount = 22
const timeOut = time.Second * 5

func Example() {
	packetConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4(127, 0, 0, 1),
		Port: 80,
	})
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, packetConn.LocalAddr().(*net.UDPAddr))
	if err != nil {
		log.Fatal(err)
	}
	lossyConn  := NewConn(conn, 1024, time.Second, time.Second*2, 0.25, UDPv4MinHeaderOverhead)
	var bytesWritten int
	go func() {
		for i := 0; i < packetCount; i++ {
			packet := make([]byte, i*64)
			bytesWritten += len(packet)
			rand.Read(packet)
			_, err := lossyConn.Write(packet)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Sent", packetCount, "packets with total size of", bytesWritten, "bytes")
	}()
	var bytesRead int
	var packetsReceived int
	startTime := time.Now()
	for {
		err = packetConn.SetDeadline(time.Now().Add(timeOut))
		if err != nil {
			log.Fatal(err)
		}
		buffer := make([]byte, 32*64)
		n, addr, err := packetConn.ReadFrom(buffer)
		if err != nil {
			if nerr, ok := err.(net.Error); ok {
				if nerr.Timeout() {
					break
				}
			}
			log.Fatal(err)
		}
		if addr.String() != conn.LocalAddr().String() {
			log.Fatal("hijacked by", addr.String())
		}
		packetsReceived++
		bytesRead += n // Ignoring the packet headers
	}
	dur := time.Now().Sub(startTime) - timeOut
	fmt.Println("Received", packetsReceived, "packets with total size of", bytesRead, "bytes in", dur.Seconds(), "seconds")
}
