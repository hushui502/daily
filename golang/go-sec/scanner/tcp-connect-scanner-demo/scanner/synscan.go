package scanner

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"time"
)


// Get the local ip and port based on our destination ip
func localIPPort(dstIP net.IP) (net.IP, int, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstIP.String()+":53421")
	if err != nil {
		return nil, 0, err
	}

	if conn, err := net.DialUDP("udp", nil, serverAddr); err != nil {
		if udpAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
			return udpAddr.IP, udpAddr.Port, nil
		}
	}

	return nil, -1, err
}

func SynScan(dstIP string, dstPort int) (string, int, error) {
	// local ip and port
	srcIP, srcPort, err := localIPPort(net.ParseIP(dstIP))
	// 可以访问的dst ip addrs
	dstAddrs, err := net.LookupIP(dstIP)
	if err != nil {
		return dstIP, 0, err
	}

	// 一个ipv4的dst ip addr
	dstip := dstAddrs[0].To4()
	// 生成类型合适的dst 和 src 端口
	var dstport layers.TCPPort
	dstport = layers.TCPPort(dstPort)
	srcport := layers.TCPPort(srcPort)

	// ip header for TCP checksum
	ip := &layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstip,
		Protocol: layers.IPProtocolTCP,
	}

	tcp := &layers.TCP{
		SrcPort: srcport,
		DstPort: dstport,
		SYN:     true,
	}

	err = tcp.SetNetworkLayerForChecksum(ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		return dstIP, 0, err
	}

	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		return dstIP, 0, err
	}
	defer conn.Close()

	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstip}); err != nil {
		return dstIP, 0, err
	}

	// set deadline so we do not wait forever.
	if err := conn.SetDeadline(time.Now().Add(4 * time.Second)); err != nil {
		return dstIP, 0, err
	}

	for {
		b := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			return dstIP, 0, err
		} else if addr.String() == dstip.String() {
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				if tcp.DstPort == srcport {
					if tcp.SYN && tcp.ACK {
						return dstIP, dstPort, err
					} else {
						return dstIP, 0, err
					}
				}
			}
		}
	}

}