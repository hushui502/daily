package scanner

import (
	"fmt"
	"net"
	"time"
)

// 建立一个tcp链接
func Connect(ip string, port int) (string, int, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 1*time.Second)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	return ip, port, err
}
