package scanner

import (
	"fmt"
	"net"
	"testing"
)

func TestLocalPort(t *testing.T) {
	var dstIP net.IP
	dstIP = []byte("127.0.0.1")
	localIp, localPort, err := localIPPort(dstIP)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(localIp.String(), "===>", localPort)
}
