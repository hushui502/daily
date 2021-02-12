package net

// eth0 ...
type HardwareType uint16

const (
	HardwareTypeLoopback = 0x0000
	HardwareTypeEthernet = 0x0001
)

func (t HardwareType) String() string {
	switch t {
	case HardwareTypeEthernet:
		return "Ethernet"
	case HardwareTypeLoopback:
		return "Loopback"
	default:
		return "Unknown"
	}
}

type EthernetType uint16

const (
	EthernetTypeIP   EthernetType = 0x0800
	EthernetTypeARP  EthernetType = 0x0806
	EthernetTypeIPv6 EthernetType = 0x86dd
)

func (t EthernetType) String() string {
	switch t {
	case EthernetTypeARP:
		return "ARP"
	case EthernetTypeIP:
		return "IP"
	case EthernetTypeIPv6:
		return "IPv6"
	default:
		return "Unknown"
	}
}

type ProtocolNumber uint8

const (
	ProtocolNumberICMP ProtocolNumber = 0x01
	ProtocolNumberTCP  ProtocolNumber = 0x06
	ProtocolNumberUDP  ProtocolNumber = 0x11
)

func (t ProtocolNumber) String() string {
	switch t {
	case ProtocolNumberICMP:
		return "ICMP"
	case ProtocolNumberTCP:
		return "TCP"
	case ProtocolNumberUDP:
		return "UDP"
	default:
		return "Unknown"
	}
}
