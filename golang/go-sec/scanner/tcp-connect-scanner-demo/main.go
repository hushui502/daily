package main

import (
	"fmt"
	"go-sec/scanner/tcp-connect-scanner-demo/scanner"
	"go-sec/scanner/tcp-connect-scanner-demo/util"
	"os"
)

// go run main.go "192.168.0.1,223.112.423.51" "12-34,1-38"
func main() {
	if len(os.Args) == 3 {
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		if err != nil {
			fmt.Println(err)
		}
		for _, ip := range ips {
			for _, port := range ports {
				fmt.Printf("ip => %s, port => %d\n", ip.String(), port)
				_, err := scanner.Connect(ip.String(), port)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("ip: %v, port: %v is open\n", ip, port)
			}
		}
	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}
