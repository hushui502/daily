package ip

import "lectcp/pkg/net"

func init()  {
	repo = newRouteTable()
	net.RegisterProtocol(net.EthernetTypeIP, rxHandler)
}

func rxHandler(dev *net.Device, data []byte, src, dst net.HardwareAddress) error {
	dgram, err := parse(data)
	if err != nil {
		return err
	}

	for _, one := range dev.Interfaces() {
		iface, ok := one.(*Interface)
		if !ok {
			continue
		}
		if iface.unicast != dgram.Dst {
			if iface.broadcast != dgram.Dst && BroadcastAddress != dgram.Dst {
				continue
			}
		}

		entry, ok := protocols[net.ProtocolNumber(dgram.Protocol)]
		if ok {
			if err := entry.rxHandler(iface, dgram.payload, dgram.Src, dgram.Dst); err != nil {
				return err
			}
		}
	}

	return nil
}
