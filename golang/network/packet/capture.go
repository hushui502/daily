package main

import (
	"fmt"
	"os"
	"time"
    "log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

var (
	deviceName  string = "eth0"
	snapshotLen int32  = 1024
	promiscuous bool   = false
	err         error
	timeout     time.Duration = -1 * time.Second
	handle      *pcap.Handle
	packetCount int = 0
)

func main() {
	// Open output pcap file and write header
	f, _ := os.Create("test.pcap")
	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
	defer f.Close()

	// Open the device for capturing
	handle, err = pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		fmt.Printf("Error opening device %s: %v", deviceName, err)
		os.Exit(1)
	}
	defer handle.Close()

	// Set filter
	var filter string = "tcp and port 80"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	// Start processing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packetCount++

		if packetCount > 100 {
			break
		}
	}
}
