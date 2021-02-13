package udp

import (
	"fmt"
	"tinytcp/ip"
	"tinytcp/net"
)

func init() {
	repo = newCbRepository()
	ip.RegisterProtocol(net.ProtocolNumberUDP, rxHandler)
}

func rxHandler(iface net.ProtocolInterface, data []byte, src, dst net.ProtocolAddress) error {
	datagram, err := parse(data, src, dst)
	if err != nil {
		return err
	}
	addr := &Address{
		Addr: iface.Address(),
		Port: datagram.DestinationPort,
	}
	entry := repo.lookup(addr)
	if entry == nil {
		return fmt.Errorf("port unreachable")
	}
	queueEntry := &queueEntry{
		addr: src,
		port: datagram.SourcePort,
		data: datagram.data,
	}
	select {
	case entry.rxQueue <- queueEntry:
		return nil
	default:
		return fmt.Errorf("drop")
	}
}
