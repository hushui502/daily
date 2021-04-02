package sensor

import (
	"encoding/json"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"go-sec/analysis/sensor/misc"
	"go-sec/analysis/sensor/models"
	"go-sec/analysis/sensor/settings"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func processPacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		if ip != nil {
			switch ip.Protocol {
			// only TCP
			case layers.IPProtocolTCP:
				tcpLayer := packet.Layer(layers.LayerTypeTCP)
				if tcpLayer != nil {
					tcp, _ := tcpLayer.(*layers.TCP)
					srcIp := ip.SrcIP.String()
					dstIp := ip.DstIP.String()
					srcPort := tcp.SrcPort.String()
					dstPort := tcp.DstPort.String()
					connInfo := models.NewConnectionInfo("tcp", srcIp, srcPort, dstIp, dstPort)

					go func(u string, info *models.ConnectionInfo) {
						if tcp.SYN && !tcp.ACK && !CheckSelfPacker(ApiUrl, connInfo) {
							misc.Log.Debugf("[TCP] %v:%v -> %v:%v, syn: %v, ack: %v, seq: %v, ack: %v, psh: %v", srcIp, srcPort,
								dstIp, dstPort, tcp.SYN, tcp.ACK, tcp.Seq, tcp.Ack, tcp.PSH)
							_ = SendPacker(info)
						}
					}(ApiUrl, connInfo)
				}
			}
		}
	}
}

func parseDNS(packet gopacket.Packet) {
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var udp layers.UDP
	var dns layers.DNS

	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet, &eth, &ip4, &udp, &dns)
	decodedLayers := make([]gopacket.LayerType, 0)
	err := parser.DecodeLayers(packet.Data(), &decodedLayers)
	if err != nil {
		return
	}
	srcIp := ip4.SrcIP
	dstIp := ip4.DstIP
	for _, v := range dns.Questions {
		dns := models.NewDns(srcIp.String(), dstIp.String(), v.Type.String(), string(v.Name))
		go func(dns *models.Dns) {
			misc.Log.Debugf("%v -> %v, dns_type: %v, dns_name: %v", srcIp, dstIp, v.Type, string(v.Name))
			_ = SendDns(dns)
		}(dns)
	}
}

func SendPacker(connInfo *models.ConnectionInfo) error {
	infoJson, err := json.Marshal(connInfo)
	if err != nil {
		return err
	}
	timestamp := time.Now().Format("2004-01-12 12:23:43")
	urlApi := fmt.Sprintf("%v%v", ApiUrl, "/api/packet/")
	secureKey := misc.MakeSign(timestamp, SecureKey)

	resp, err := http.PostForm(urlApi, url.Values{"timestamp": {timestamp}, "secureKey": {secureKey}, "data": {string(infoJson)}})
	if err != nil {
		return err
	}
	// TODO: resp.Status
	_ = resp

	return nil
}

func SendDns(dns *models.Dns) error {
	reqJson, err := json.Marshal(dns)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	urlApi := fmt.Sprintf("%v%v", ApiUrl, "/api/dns/")
	secureKey := misc.MakeSign(timestamp, SecureKey)
	_, err = http.PostForm(urlApi, url.Values{"timestamp": {timestamp}, "secureKey": {secureKey}, "data": {string(reqJson)}})
	return err
}

func CheckSelfPacker(apiUrl string, p *models.ConnectionInfo) bool {
	urlParsed, err := url.Parse(apiUrl)
	if err != nil {
		return false
	}
	apiHost := urlParsed.Host
	apiIp := strings.Split(apiHost, ":")[0]
	sensorIp := settings.Ips[0]

	if p.SrcIp == sensorIp && p.DstIp == apiIp || p.SrcIp == apiIp && p.DstIp == sensorIp {
		return true
	}

	return false
}
