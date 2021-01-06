package main

import (
	"fmt"
	"net"
	"os"
)

// hosts may have multiple IP address, usually from multiple network interface cards.
// They may also have multiple host names, acting as aliases.

func main() {
	if len(os.Args) != 2 {
		return
	}

	name := os.Args[1]

	addrs, err := net.LookupHost(name)
	if err != nil {
		return
	}

	for _, addr := range addrs {
		fmt.Println(addr)
	}

	cname, _ := net.LookupCNAME(name)
	fmt.Println(cname)


	os.Exit(1)


}
