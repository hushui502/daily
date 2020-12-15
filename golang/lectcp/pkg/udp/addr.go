package udp

import (
	"fmt"
	"lectcp/pkg/net"
)

type Address struct {
	Addr net.ProtocolAddress
	Port uint16
}

func (a Address) Newwork() string {
	return "udp"
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Addr, a.Port)
}
