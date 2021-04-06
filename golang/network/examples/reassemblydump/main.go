package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

var (
	maxcount         = flag.Int("c", -1, "Only grab this many packets, then exit")
	decoder          = flag.String("decoder", "", "Name of the decoder to use (default: guess from capture)")
	statsevery       = flag.Int("stats", 1000, "Output statistics every N packets")
	lazy             = flag.Bool("lazy", false, "If true, do lazy decoding")
	nodefrag         = flag.Bool("nodefrag", false, "If true, do not do IPv4 defrag")
	checksum         = flag.Bool("checksum", false, "Check TCP checksum")
	nooptcheck       = flag.Bool("nooptcheck", false, "Do not check TCP options (useful to ignore MSS on captures with TSO)")
	ignorefsmerr     = flag.Bool("ignorefsmerr", false, "Ignore TCP FSM errors")
	allowmissinginit = flag.Bool("allowmissinginit", false, "Support streams without SYN/SYN+ACK/ACK sequence")
	verbose          = flag.Bool("verbose", false, "Be verbose")
	debug            = flag.Bool("debug", false, "Display debug information")
	quiet            = flag.Bool("quiet", false, "Be quiet regarding errors")
)

// http
var (
	nohttp          = flag.Bool("nohttp", false, "Disable HTTP parsing")
	output          = flag.String("output", "", "Path to create file for HTTP 200 OK responses")
	writeincomplete = flag.Bool("writeincomplete", false, "Write incomplete response")
	hexdump         = flag.Bool("dump", false, "Dump HTTP request/response as hex")
	hexdumppkt      = flag.Bool("dumppkt", false, "Dump packet as hex")
)

// captures
var (
	iface   = flag.String("i", "eth0", "Interface to read packets from")
	fname   = flag.String("r", "", "Filename to read from, overrides -i")
	snaplen = flag.Int("s", 65536, "Snap length (number of bytes max to read per packet")
	tstype  = flag.String("timestamp_type", "", "Type of timestamps to use")
	promisc = flag.Bool("promisc", true, "Set promiscuous mode")
)

// profile
var (
	memprofile = flag.String("memprofile", "", "Write memory profile")
)

// FSM
var stats struct {
	ipdefrag            int
	missedBytes         int
	pkt                 int
	sz                  int
	totalsz             int
	rejectFsm           int
	rejectOpt           int
	rejectConnFsm       int
	reassembled         int
	outOfOrderBytes     int
	outOfOrderPackets   int
	biggestChunkBytes   int
	biggestChunkPackets int
	overlapBytes        int
	overlapPackets      int
}

const closeTimeout time.Duration = time.Hour * 24
const timeout time.Duration = time.Minute * 5

/*
	HTTP part
*/
type httpReader struct {
	ident    string
	isClient bool
	bytes    chan []byte
	data     []byte
	hexdump  bool
	parent   *tcpStream
}

func (h *httpReader) Read(p []byte) (int, error) {
	ok := true
	for ok && len(h.data) == 0 {
		h.data, ok = <-h.bytes
	}
	if !ok || len(h.data) == 0 {
		return 0, io.EOF
	}
	l := copy(p, h.data)
	return l, nil
}

// error process
var (
	outputLevel    int
	errorsMap      map[string]int
	errorsMapMutex sync.Mutex
	errors         uint
)

func Error(t string, s string, a ...interface{}) {
	errorsMapMutex.Lock()
	errors++
	nb, _ := errorsMap[t]
	errorsMap[t] = nb + 1
	errorsMapMutex.Unlock()
	if outputLevel >= 0 {
		fmt.Printf(s, a...)
	}
}

func Info(s string, a ...interface{}) {
	if outputLevel >= 1 {
		fmt.Printf(s, a...)
	}
}

func Debug(s string, a ...interface{}) {
	if outputLevel >= 2 {
		fmt.Printf(s, a...)
	}
}

func (h *httpReader) run(wg *sync.WaitGroup) {
	defer wg.Done()

	b := bufio.NewReader(h)
	for {
		if h.isClient {
			req, err := http.ReadRequest(b)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				Error("HTTP-request", "HTTP/%s Request error: %s (%v,%+v)\n", h.ident, err, err, err)
				continue
			}
			body, err := ioutil.ReadAll(req.Body)
			s := len(body)
			if err != nil {
				Error("HTTP-request-body", "Got body err: %s\n", err)
			} else if h.hexdump {
				Info("Body(%d/0x%x)\n%s\n", len(body), len(body), hex.Dump(body))
			}
			req.Body.Close()
			Info("HTTP/%s Request: %s %s (body:%d)\n", h.ident, req.Method, req.URL, s)
			h.parent.Lock()
			h.parent.urls = append(h.parent.urls, req.URL.String())
			h.parent.Unlock()
		} else {
			res, err := http.ReadResponse(b, nil)
			var req string
			h.parent.Lock()
			if len(h.parent.urls) == 0 {
				req = fmt.Sprintf("<no-request-seen>")
			} else {
				req, h.parent.urls = h.parent.urls[0], h.parent.urls[1:]
			}
			h.parent.Unlock()
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				Error("HTTP-response", "HTTP/%s Response error: %s (%v,%+v)\n", h.ident, err, err, err)
				continue
			}
			body, err := ioutil.ReadAll(res.Body)
			s := len(body)
			if err != nil {
				Error("HTTP-response-body", "HTTP/%s: failed to get body(parsed len:%d): %s\n", h.ident, s, err)
			}
			if h.hexdump {
				Info("Body(%d/0x%x)\n%s\n", len(body), len(body), hex.Dump(body))
			}
			res.Body.Close()
			sym := ","
			if res.ContentLength > 0 && res.ContentLength != int64(s) {
				sym = "!="
			}
			contentType, ok := res.Header["Content-Type"]
			if !ok {
				contentType = []string{http.DetectContentType(body)}
			}
			encoding := res.Header["Content-Encoding"]
			Info("HTTP/%s Response: %s URL:%s (%d%s%d%s) -> %s\n", h.ident, res.Status, req, res.ContentLength, sym, s, contentType, encoding)
			if (err == nil || *writeincomplete) && *output != "" {
				base := url.QueryEscape(path.Base(req))
				if err != nil {
					base = "incomplete-" + base
				}
				base = path.Join(*output, base)
				if len(base) > 250 {
					base = base[:250] + "..."
				}
				if base == *output {
					base = path.Join(*output, "noname")
				}
				target := base
				n := 0
				for {
					_, err := os.Stat(target)
					if err != nil {
						break
					}
					target = fmt.Sprintf("%s-%s", base, n)
					n++
				}
				f, err := os.Create(target)
				if err != nil {
					Error("HTTP-create", "Cannot create %s: %s\n", target, err)
					continue
				}
				var r io.Reader
				r = bytes.NewBuffer(body)
				if len(encoding) > 0 && (encoding[0] == "gzip" || encoding[0] == "deflate") {
					r, err = gzip.NewReader(r)
					if err != nil {
						Error("HTTP-gunzip", "Failed to gzip decode: %s", err)
					}
				}
				if err != nil {
					w, err := io.Copy(f, r)
					if _, ok := r.(*gzip.Reader); ok {
						r.(*gzip.Reader).Close()
					}
					f.Close()
					if err != nil {
						Error("HTTP-save", "%s: failed to save %s (l:%d): %s\n", h.ident, target, w, err)
					} else {
						Info("%s: Saved %s (l:%d)\n", h.ident, target, w)
					}
				}
			}
		}
	}
}

/*
	TCP factory: returns a new Stream
*/
type tcpStreamFactory struct {
	wg     sync.WaitGroup
	doHTTP bool
}

func (f *tcpStreamFactory) New(net, transport gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	Debug("* NEW: %s %s\n", net, transport)

	fsmOption := reassembly.TCPSimpleFSMOptions{
		SupportMissingEstablishment: *allowmissinginit,
	}
	stream := &tcpStream{
		net:        net,
		transport:  transport,
		isDNS:      tcp.SrcPort == 53 || tcp.DstPort == 53,
		isHTTP:     (tcp.SrcPort == 80 || tcp.DstPort == 80) && f.doHTTP,
		reversed:   tcp.SrcPort == 80,
		tcpstate:   reassembly.NewTCPSimpleFSM(fsmOption),
		ident:      fmt.Sprintf("%s:%s", net, transport),
		optchecker: reassembly.NewTCPOptionCheck(),
	}
	if stream.isHTTP {
		stream.client = httpReader{
			bytes:    make(chan []byte),
			ident:    fmt.Sprintf("%s %s", net, transport),
			hexdump:  *hexdump,
			parent:   stream,
			isClient: true,
		}
		stream.server = httpReader{
			bytes:   make(chan []byte),
			ident:   fmt.Sprintf("%s %s", net.Reverse(), transport.Reverse()),
			hexdump: *hexdump,
			parent:  stream,
		}
		f.wg.Add(2)
		go stream.client.run(&f.wg)
		go stream.server.run(&f.wg)
	}
	return stream
}

func (f *tcpStreamFactory) WaitGoroutines() {
	f.wg.Wait()
}

/*
	assembler context
*/
type Context struct {
	CaptureInfo gopacket.CaptureInfo
}

func (c *Context) GetCaptureInfo() gopacket.CaptureInfo {
	return c.CaptureInfo
}

/*
	TCP stream
*/
/* It's a connection (bidirectional) */
type tcpStream struct {
	tcpstate       *reassembly.TCPSimpleFSM
	fsmerr         bool
	optchecker     reassembly.TCPOptionCheck
	net, transport gopacket.Flow
	isDNS          bool
	isHTTP         bool
	reversed       bool
	client         httpReader
	server         httpReader
	urls           []string
	ident          string
	sync.Mutex
}

func (t *tcpStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, netxSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	// FSM final-state machine
	if !t.tcpstate.CheckState(tcp, dir) {
		Error("FSM", "%s: Packet rejected by FSM (state:%s)\n", t.ident, t.tcpstate.String())
		stats.rejectFsm++
		if !t.fsmerr {
			t.fsmerr = true
			stats.rejectConnFsm++
		}
		if !*ignorefsmerr {
			return false
		}
	}
	// Options
	err := t.optchecker.Accept(tcp, ci, dir, netxSeq, start)
	if err != nil {
		Error("OptionChecker", "%s: Packet rejected by OptionChecker: %s\n", t.ident, err)
		stats.rejectOpt++
		if !*nooptcheck {
			return false
		}
	}
	// Checksum
	accept := true
	if *checksum {
		c, err := tcp.ComputeChecksum()
		if err != nil {
			Error("ChecksumCompute", "%s: Got error computing checksum: %s\n", t.ident, err)
			accept = false
		} else if c != 0x0 {
			Error("Checksum", "%s: Invalid checksum: 0x%x\n", t.ident, c)
			accept = false
		}
	}
	if !accept {
		stats.rejectOpt++
	}

	return accept
}

func (t *tcpStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, start, end, skip := sg.Info()
	length, saved := sg.Lengths()

	sgStats := sg.Stats()
	if skip > 0 {
		stats.missedBytes += skip
	}
	stats.sz += length - saved
	stats.pkt += sgStats.Packets
	if sgStats.Chunks > 1 {
		stats.reassembled++
	}
	stats.outOfOrderBytes += sgStats.QueuedBytes
	stats.outOfOrderPackets += sgStats.QueuedPackets
	if length > stats.biggestChunkBytes {
		stats.biggestChunkBytes = length
	}
	if sgStats.Packets > stats.biggestChunkPackets {
		stats.biggestChunkPackets = length
	}
	if sgStats.OverlapBytes != 0 && sgStats.OverlapPackets == 0 {
		fmt.Printf("bytes:%d, pkts:%d\n", sgStats.OverlapBytes, sgStats.OverlapPackets)
		panic("Invalid overlap")
	}
	stats.overlapBytes += sgStats.OverlapBytes
	stats.overlapPackets += sgStats.OverlapPackets

	var ident string
	if dir == reassembly.TCPDirClientToServer {
		ident = fmt.Sprintf("%v %v(%s): ", t.net, t.transport, dir)
	} else {
		ident = fmt.Sprintf("%v %v(%s): ", t.net.Reverse(), t.transport.Reverse(), dir)
	}
	Debug("%s: SG reassembled packet with %d bytes (start:%v,end:%v,skip:%d,saved:%d,nb:%d,%d,overlap:%d,%d)\n",
		ident, length, start, end, skip, saved, sgStats.Packets, sgStats.Chunks, sgStats.OverlapBytes, sgStats.OverlapPackets)
	if skip == -1 && *allowmissinginit {
		// pass
	} else if skip != 0 {
		// missing bytes in stream: do not even try to parse it
		return
	}

	data := sg.Fetch(length)
	if t.isDNS {
		dns := &layers.DNS{}
		var decoded []gopacket.LayerType
		if len(data) < 2 {
			if len(data) > 0 {
				sg.KeepFrom(0)
			}
			return
		}
		dnsSize := binary.BigEndian.Uint16(data[:2])
		missing := int(dnsSize) - len(data[2:])
		Debug("dnsSize: %d, missing: %d\n", dnsSize, missing)
		if missing > 0 {
			Info("Missing some bytes: %d\n", missing)
			sg.KeepFrom(0)
			return
		}
		p := gopacket.NewDecodingLayerParser(layers.LayerTypeDNS, dns)
		err := p.DecodeLayers(data[2:], &decoded)
		if err != nil {
			Error("DNS-parser", "Failed to decode DNS: %v\n", err)
		} else {
			Debug("DNS: %s\n", gopacket.LayerDump(dns))
		}
		if len(data) > 2+int(dnsSize) {
			sg.KeepFrom(2 + int(dnsSize))
		}
	} else if t.isHTTP {
		if length > 0 {
			if *hexdump {
				Debug("Feeding http with:\n%s", hex.Dump(data))
			}
			if dir == reassembly.TCPDirClientToServer && !t.reversed {
				t.client.bytes <- data
			} else {
				t.server.bytes <- data
			}
		}
	}
}

func (t *tcpStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	Debug("%s: Connection closed\n", t.ident)
	if t.isHTTP {
		close(t.client.bytes)
		close(t.server.bytes)
	}
	// don't remove connection for the last ACK
	return false
}
