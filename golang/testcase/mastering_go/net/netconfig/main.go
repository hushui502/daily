package main

import (
	"fmt"
	"net"
)

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range interfaces {
		fmt.Printf("Interface: %v\n", i.Name)
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			fmt.Println(err)
		}
		address, _ := byName.Addrs()
		for k, v := range address {
			fmt.Println("Intserface address: ", k, v.String())
		}
	}

}
