package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {
	fmt.Println("-----find all devices----")

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Println("======DEVICE NAME: ", device.Name)
		fmt.Println("======DEVICE Description: ", device.Description)
		fmt.Println("======DEVICE Flag: ", device.Flags)

		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
			fmt.Println("- Board : ", address.Broadaddr)
			fmt.Println("- P2P : ", address.P2P)
			fmt.Println("=====================")
		}
	}
}
