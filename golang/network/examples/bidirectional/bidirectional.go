package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"log"
	"network/examples/util"
	"time"
)

var iface = flag.String("i", "eth0", "Interface to get packets from")
var snaplen = flag.Int("s", 16<<10, "Snaplen for pcap packet capture")
var filter = flag.String("f", "tcp", "BFT filter for pacp")
var logAllPackets = flag.Bool("v", false, "Logs every packet in greate detail")

// key is used to map bidirectional streams to each other
type key struct {
	net, transport gopacket.Flow
}

// rewrite String for a human-readable fashion
func (k key) String() string {
	return fmt.Sprintf("%v:%v", k.net, k.transport)
}

const timeout time.Duration = time.Minute * 5

type myStream struct {
	bytes int64
	bidi  *bidi
	done  bool
}

type bidi struct {
	key            key
	a, b           *myStream
	lastPacketSeen time.Time
}

type myFactory struct {
	// bidiMap maps keys to bidirectional stream pairs.
	bidiMap map[key]*bidi
}

func (f *myFactory) New(netFlow, tcpFlow gopacket.Flow) tcpassembly.Stream {
	s := &myStream{}

	k := key{netFlow, tcpFlow}
	bd := f.bidiMap[k]
	if bd == nil {
		bd = &bidi{a: s, key: k}
		log.Printf("[%v] created first side of bidirectional stream", bd.key)
		// Register bidirectional with the reverse key, so the matching stream going
		// the other direction will find it. -.-
		f.bidiMap[key{netFlow.Reverse(), tcpFlow.Reverse()}] = bd
	} else {
		log.Printf("[%v] found second side of bidirectional stream", bd.key)
		bd.b = s
		delete(f.bidiMap, k)
	}
	s.bidi = bd

	return s
}

// emptyStream is used to finished bidi that only have one stream in collectOldStreams
// and sets/finishes the 'b' stream inside them.
var emptyStream = &myStream{done: true}

// collectOldStreams finds any streams that haven't received a packet within timeout
func (f *myFactory) collectOldStreams() {
	cutOff := time.Now().Add(-timeout)
	for k, bd := range f.bidiMap {
		if bd.lastPacketSeen.Before(cutOff) {
			log.Printf("[%v] timing out old stream", bd.key)
			bd.b = emptyStream
			delete(f.bidiMap, k)
			bd.maybeFinish() // if b was the last stream we were waiting for, finish up.
		}
	}
}

// reassembled handles reassembled TCP stream data
func (s *myStream) Reassembled(rs []tcpassembly.Reassembly) {
	for _, r := range rs {
		s.bytes += int64(len(r.Bytes))
		if r.Skip > 0 {
			s.bytes += int64(r.Skip)
		}
		// Mark that we've received new packet data.
		if s.bidi.lastPacketSeen.Before(r.Seen) {
			s.bidi.lastPacketSeen = r.Seen
		}
	}
}

// ReassemblyComplete marks this stream as finished.
func (s *myStream) ReassemblyComplete() {
	s.done = true
	s.bidi.maybeFinish()
}

func (bd *bidi) maybeFinish() {
	switch {
	case bd.a == nil:
		log.Fatalf("[%v] a should always be non-nil, since it's set when bidis are created", bd.key)
	case !bd.a.done:
		log.Printf("[%v] still waiting on first stream", bd.key)
	case bd.b == nil:
		log.Printf("[%v] no second stream yet", bd.key)
	case !bd.b.done:
		log.Printf("[%v] still waiting on second stream", bd.key)
	default:
		log.Printf("[%v] FINISHED, bytes: %d tx, %d rx", bd.key, bd.a.bytes, bd.b.bytes)
	}
}

func main() {
	defer util.Run()()
	log.Printf("starting capture on interface %q", *iface)

	handle, err := pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}
	if err := handle.SetBPFFilter(*filter); err != nil {
		panic(err)
	}

	// setup assembler  (template)
	streamFactory := &myFactory{bidiMap: make(map[key]*bidi)}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	// limit memory usage by auto-flushing connection state if we get over 100k packets in memory.
	assembler.MaxBufferedPagesTotal = 100000
	assembler.MaxBufferedPagesPerConnection = 1000

	log.Println("reading in packets")

	// read in packets, pass to assembler
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(timeout / 4)

	for {
		select {
		case packet := <-packets:
			if *logAllPackets {
				log.Println(packets)
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			// flush connections that we have not seen activity in the past minute
			log.Println("---- FLUSHING ----")
			assembler.FlushOlderThan(time.Now().Add(-timeout))
			streamFactory.collectOldStreams()

		}
	}
}
