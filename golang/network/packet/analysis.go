package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"strings"
	"time"
)

//var (
//	pcapFile string = "D:\\project\\go\\src\\lib\\practice\\high-performance-go\\netprotocol\\ssl\\test.pcap"
//	handle *pcap.Handle
//	err error
//)
//
//func main() {
//	// Open file instead of device
//	handle, err = pcap.OpenOffline(pcapFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer handle.Close()
//
//	// Loop through packets in file
//	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
//	for packet := range packetSource.Packets() {
//		fmt.Println(packet)
//	}
//}

var (
	device      string = "eth0"
	snapshotLen int32  = 1024
	promiscuous bool   = false
	err         error
	timeout     time.Duration = 30 * time.Second
	handle      *pcap.Handle
)

func main() {
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {log.Fatal(err) }
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacketInfo(packet)
	}
}

func printPacketInfo(packet gopacket.Packet) {

	// ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Println("Ethernet layer detected")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("Source Mac: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination Mac: ", ethernetPacket.DstMAC)
		fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		fmt.Println()
	}

	// ip packet
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (4 or 6)
		// THL (IP Header length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol, Checksum. SrcIP, DstIP
		fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
		fmt.Println()
	}

	// tcp packet
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)

		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool Flags: FIN, SYN, RST, RSH, ACK, URG, ECE, CWR, NS
		fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Println("Sequence number: ", tcp.Seq)
		fmt.Println()
	}

	// Iterate over all layers, printing out each layer type
	fmt.Println("All packet layers:")
	for _, layer := range packet.Layers() {
		fmt.Println("- ", layer.LayerType())
	}

	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		fmt.Printf("%s\n", applicationLayer.Payload())

		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			fmt.Println("HTTP found!")
		}
	}

	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}

}