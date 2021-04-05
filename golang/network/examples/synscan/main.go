package main

import (
	"errors"
	"flag"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/routing"
	"log"
	"net"
	"network/examples/util"
	"time"
)

// scanner handles scanning a single IP address
type scanner struct {
	// iface is the interface to send packets on.
	iface *net.Interface
	// destination, gateway(if applicable), and source IP address to use.
	dst, gw, src net.IP

	handle *pcap.Handle

	// opts and buf allow us to easily serialize packets in the address,
	// using router to determine how to route pakets to that IP.
	opts gopacket.SerializeOptions
	buf gopacket.SerializeBuffer
}

func newScanner(ip net.IP, router routing.Router) (*scanner, error) {
	s := &scanner{
		dst: ip,
		opts: gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
		buf: gopacket.NewSerializeBuffer(),
	}
	// figure out the route to the IP
	iface, gw, src, err := router.Route(ip)
	if err != nil {
		return nil, err
	}
	log.Printf("scanning ip %v with interface %v, gateway %v, src %v", ip, iface.Name, gw, src)
	s.gw, s.src, s.iface = gw, src, iface

	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	s.handle = handle
	return s, nil
}

func (s *scanner) close() {
	s.handle.Close()
}

func (s *scanner) getHwAddr() (net.HardwareAddr, error) {
	start := time.Now()
	arpDst := s.dst
	if s.gw != nil {
		arpDst = s.gw
	}
	// prepare the layers to send for an ARP request
	eth := layers.Ethernet{
		SrcMAC: s.iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:        layers.LinkTypeEthernet,
		Protocol:        layers.EthernetTypeIPv4,
		HwAddressSize:   6,
		ProtAddressSize: 4,
		Operation:       layers.ARPRequest,
		SourceHwAddress: []byte(s.iface.HardwareAddr),
		SourceProtAddress: []byte(s.src),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(arpDst),
	}

	// send a single arp request packet
	if err := s.send(&eth, &arp); err != nil {
		return nil, err
	}

	// wait 3 seconds for an ARP reply
	for {
		if time.Since(start) > time.Second * 3 {
			return nil, errors.New("timeout getting ARP reply")
		}
		data, _, err := s.handle.ReadPacketData()
		if err == pcap.NextErrorTimeoutExpired {
			continue
		} else if err != nil {
			return nil, err
		}
		// dst -> src
		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
			arp := arpLayer.(*layers.ARP)
			if net.IP(arp.SourceProtAddress).Equal(net.IP(arpDst)) {
				return net.HardwareAddr(arp.SourceHwAddress), nil
			}
		}
	}
}

// scan scans the dst IP address of this scanner
func (s *scanner) scan() error {
	hwaddr, err := s.getHwAddr()
	if err != nil {
		return err
	}
	// construct all the network layers we need
	eth := layers.Ethernet{
		SrcMAC: s.iface.HardwareAddr,
		DstMAC: hwaddr,
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip4 := layers.IPv4{
		SrcIP:   s.src,
		DstIP:   s.dst,
		Version: 4,
		TTL:     64,
		Protocol: layers.IPProtocolTCP,
	}
	tcp := layers.TCP{
		SrcPort: 54321,
		DstPort: 0,
		SYN: true,
	}
	tcp.SetNetworkLayerForChecksum(&ip4)

	// create the flow we expect returning packets to have,
	// so we can check against it and discard useless packets.
	ipFlow := gopacket.NewFlow(layers.EndpointIPv4, s.dst, s.src)
	start := time.Now()
	for {
		// send one packet peer loop iteration until we have sent packets
		// to all of ports[1, 65535]
		if tcp.DstPort < 65535 {
			start = time.Now()
			tcp.DstPort++
			if err := s.send(&eth, &ip4, &tcp); err != nil {
				log.Printf("error sending to port %v: %v", tcp.DstPort, err)
			}
		}
		// time out 5 seconds after the last packet we sent.
		if time.Since(start) > time.Second * 5 {
			log.Printf("timed out for %v, assuming we've seen all we can", s.dst)
			return nil
		}

		// read in the next packet.
		data, _, err := s.handle.ReadPacketData()
		if err == pcap.NextErrorTimeoutExpired {
			continue
		} else if err != nil {
			return err
		}

		// parse the packet. we'd use DecodingLayerParser here if we
		// want to be really fast.
		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)

		// find the packets we care about, and print out logging information
		// about them. all others are ignored
		if net := packet.NetworkLayer(); net == nil {
			//log.Printf("packet has no network layer")
		} else if net.NetworkFlow() != ipFlow {
			// log.Printf("packet does not match our ip src/dst")
		} else if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer == nil {
			// log.Printf("packet has not tcp layer")
		} else if tcp, ok := tcpLayer.(*layers.TCP); !ok {
			panic("tcp layer is not tcp layer :-/")
		} else if tcp.DstPort != 54321 {
			// log.Printf("dst port %v does not match", tcp.DstPort)
		} else if tcp.RST {
			log.Printf("  port %v closed", tcp.SrcPort)
		} else if tcp.SYN && tcp.ACK {
			log.Printf("  port %v open", tcp.SrcPort)
		} else {
			// log.Printf("ignoring useless packet")
		}
	}
}

func (s *scanner) send(l ...gopacket.SerializableLayer) error {
	if err := gopacket.SerializeLayers(s.buf, s.opts, l...); err != nil {
		return err
	}
	return s.handle.WritePacketData(s.buf.Bytes())
}

func main() {
	defer util.Run()
	router, err := routing.New()
	if err != nil {
		log.Fatal("routing error:", err)
	}

	for _, arg := range flag.Args() {
		var ip net.IP
		if ip = net.ParseIP(arg); ip == nil {
			log.Printf("non-ip target: %q", arg)
			continue
		} else if ip = ip.To4(); ip == nil {
			log.Printf("non-ipv4 target: %q", arg)
			continue
		}

		s, err := newScanner(ip, router)
		if err != nil {
			log.Printf("unable to create scanner for %v: %v", ip, err)
			continue
		}
		if err := s.scan(); err != nil {
			log.Printf("unable to scan %v: %v", ip, err)
		}
		s.close()
	}
}
