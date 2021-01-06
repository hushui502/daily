package main

import (
	"fmt"
	"net"
	"os"
)

// many of the other functions and methods in the net package return a pointer to an IPAddr. This
// is simple a structure containing an IP

// type IPAddr struct {
// 	 IP IP
// }




func main() {
	if len(os.Args) != 2 {
		// ...
		return
	}

	name := os.Args[1]
	// a primary use of this type is to perform DNS lookups on IP host names.
	// The function ResolveIPAddr will perform a DNS lookup on a hostname, and return a single IP address.
	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		// ...
		return
	}
	fmt.Println("Resolved address is ", addr.String())
	os.Exit(1)
}
// ./main www.google.com
// Resolved address is 66.102.11.104


