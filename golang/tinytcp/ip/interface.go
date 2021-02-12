package ip

import "tinytcp/net"

type Interface struct {
	unicast Address
	netmask Address
	broadcast Address
	gateway Address
	device *net.Device
}


