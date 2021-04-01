package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
	"log"
)

var (
	assembler *tcpassembly.Assembler
	decodeOptions gopacket.DecodeOptions
)

func init() {
	streamFactory := &tlsStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler = tcpassembly.NewAssembler(streamPool)
	decodeOptions = gopacket.DecodeOptions{
		SkipDecodeRecovery: true,
		DecodeStreamsAsDatagrams: true,
	}
}

type tlsStreamFactory struct {}


func (h *tlsStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	tstream := &tlsStream{
		net: net,
		transport: transport,
		r: tcpreader.NewReaderStream(),
	}

	go tstream.run()

	return &tstream.r
}

type tlsStream struct {
	net, transport gopacket.Flow
	r tcpreader.ReaderStream
}

func (t *tlsStream) run() {
	buf := bufio.NewReader(&t.r)
	for {
		rLen, err := buf.Peek(6)
		if err != nil {
			continue
		}
		len := binary.BigEndian.Uint16(rLen[3:5]) + 5 // header length
		raw := make([]byte, len)
		_, err = io.ReadFull(buf,raw)

		if err != nil {
			continue
		}

		p := gopacket.NewPacket(raw, layers.LayerTypeTLS, decodeOptions)

		fmt.Println("-------------- SSL Packet ------------")
		fmt.Println(p)
	}
}

func PacketHandler(p gopacket.Packet) {
	assembler.AssembleWithTimestamp(p.NetworkLayer().NetworkFlow(), p.TransportLayer().(*layers.TCP), p.Metadata().Timestamp)
}

func main() {
	// Open file instead of device
	handle, err := pcap.OpenOffline("ssl_test.pcap")
	if err != nil { log.Fatal(err) }
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		PacketHandler(packet)
	}
}
