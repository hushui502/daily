package udp

import (
	"fmt"
	"tinytcp/ip"
)

func Dial(local, remote *Address) (*Conn, error) {
	if local == nil {
		iface := ip.GetInterfaceByRemoteAddress(remote.Addr)
		if iface == nil {
			return nil, fmt.Errorf("dial failure")
		}
		local = &Address{
			Addr: iface.Address(),
		}
	}
	entry := repo.add(local)
	if entry == nil {
		return nil, fmt.Errorf("dial failture")
	}

	return &Conn{
		cb:   entry,
		peer: remote,
	}, nil
}

func Listen(local *Address) (*Conn, error) {
	entry := repo.add(local)
	if entry == nil {
		return nil, fmt.Errorf("listen failure")
	}

	return &Conn{cb: entry}, nil
}
