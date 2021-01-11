package util

import (
	"fmt"
	"github.com/malfunkt/iprange"
	"net"
	"strconv"
	"strings"
)

// 根据端口号字符串解析出一个合理的端口数组并返回给调用者
func GetPorts(selection string) ([]int, error) {
	ports := make([]int, 0)
	if selection == "" {
		return ports, nil
	}

	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid port selection segment: '%s'\n", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'\n", p1)
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'\n", p2)
			}

			if p1 > p2 {
				return nil, fmt.Errorf("Invalid port range: %s-%s\n", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}
		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'\n", r)
			} else {
				ports = append(ports, port)
			}
		}
	}

	return ports, nil
}

// 获取一个ip list
func GetIpList(ip string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ip)
	if err != nil {
		return nil, err
	}

	list := addressList.Expand()

	return list, err
}
