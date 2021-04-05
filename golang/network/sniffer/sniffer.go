package main

import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/sync/errgroup"
	"log"
)

type Sniffer struct {
	iface string
	ctx   context.Context
	file  string
	synCh chan string
	errgroup.Group
}

type Options func(*Sniffer)

type flowFunc func(gopacket.Flow, gopacket.Flow) string

type errGoFunc func() error

func NewSniffer(opts ...Options) *Sniffer {
	s := &Sniffer{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Sniffer) Start() error {
	lru := NewLru(100)
	packets, errCh := s.generator()

LOOP:
	for {
		select {
		case p := <-packets:
			if p == nil {
				// break for-select
				break LOOP
			}
			if ip, tcp := decode(p); ip != nil && tcp != nil {
				if !(tcp.SYN || tcp.FIN || tcp.RST || (tcp.ACK && len(tcp.LayerPayload()) == 0)) {
					b := tcp.LayerPayload()
					network, transport := ip.NetworkFlow(), tcp.TransportFlow()
					s.Go(process(lru, b, network, transport, flow, s.synCh))
				}
			}
		case <-s.ctx.Done():
			return s.ctx.Err()
		case err := <-errCh:
			if err != nil {
				log.Printf("[-] %+v\n", err)
			}
		}
	}
	return nil
}

func flow(ipflow, tcpflow gopacket.Flow) string {
	ipSrc, ipDst := ipflow.Endpoints()
	portSrc, portDst := tcpflow.Endpoints()
	return fmt.Sprintf("%s:%s %s:%s",
		ipSrc.String(), portSrc.String(), ipDst.String(), portDst.String())
}

func decode(p gopacket.Packet) (pv4 *layers.IPv4, tcp *layers.TCP) {
	if layer := p.Layer(layers.LayerTypeIPv4); layer != nil {
		ip, ok := layer.(*layers.IPv4)
		if !ok {
			return nil, nil
		}
		switch ip.Protocol {
		case layers.IPProtocolTCP:
			if l4 := p.Layer(layers.LayerTypeTCP); l4 != nil {
				if tcp, ok = l4.(*layers.TCP); !ok {
					return nil, nil
				}
				return ip, tcp
			}

		}
	}
	return nil, nil
}

func process(lru Lru, b []byte, network, transport gopacket.Flow, fl flowFunc, in chan string) errGoFunc {
	return func() error {
		if uint8(b[0]) != 22 {
			return nil
		}
		fh := network.FastHash()
		value, ok := lru.Get(fh)
		if !ok {
			flow := fl(network, transport)
			hs := NewHandShake(flow, in)
			err := hs.Unmarshal(b)
			if err != nil {
				return err
			}
			lru.Add(fh, hs)
			return nil
		}
		err := value.(*Handshake).Unmarshal(b)
		if err != nil {
			return err
		}
		lru.Add(fh, value)

		return nil
	}
}

func (s *Sniffer) Stop() error {
	if err := s.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *Sniffer) generator() (<-chan gopacket.Packet, <-chan error) {
	packet := make(chan gopacket.Packet, 1)
	errCh := make(chan error, 1)
	handler := &pcap.Handle{}
	var err error

	if len(s.file) != 0 {
		handler, err = pcap.OpenOffline(s.file)
		if err != nil {
			errCh <- err
		}
	} else {
		handler, err = pcap.OpenLive(s.iface, 65535, false, pcap.BlockForever)
		if err != nil {
			errCh <- err
		}
	}

	packetSource := gopacket.NewPacketSource(handler, handler.LinkType())

	go func() {
		defer func() {
			close(errCh)
			handler.Close()
		}()

		for {
			select {
			case <-s.ctx.Done():
				errCh <- s.ctx.Err()
				return
			case packet <- <-packetSource.Packets():
			}
		}
	}()

	return packet, errCh
}

func SetDevice(iface string) Options {
	return func(sniffer *Sniffer) {
		sniffer.iface = iface
	}
}

func SetCtx(ctx context.Context) Options {
	return func(sniffer *Sniffer) {
		sniffer.ctx = ctx
	}
}

func SetSyncCh(ch chan string) Options {
	return func(sniffer *Sniffer) {
		sniffer.synCh = ch
	}
}

func SetDump(path string) Options {
	return func(sniffer *Sniffer) {
		sniffer.file = path
	}
}
