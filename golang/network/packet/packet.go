package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
)

func main() {
	handle, err := pcap.OpenLive("eth0",
		65536, true, pcap.BlockForever)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer handle.Close()

	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()

	for {
		var packet gopacket.Packet
		select {
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			println(arp)
			if net.HardwareAddr(arp.SourceHwAddress).String() == "abc" {

			}

			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer == nil {
				continue
			}
			tcp := tcpLayer.(*layers.TCP)
			println(tcp)

			udpLayer := packet.Layer(layers.LayerTypeUDP)
			if udpLayer == nil {
				continue
			}
			udp := udpLayer.(*layers.UDP)
			println(udp)

		}
	}
}
