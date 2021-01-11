package main

import (
	"fmt"
	"go-sec/scanner/tcp-connect-scanner-demo/scanner"
	"go-sec/scanner/tcp-connect-scanner-demo/util"
	"os"
)

// go run main.go "192.168.0.1,223.112.423.51" "12-34,1-38"
//func main() {
//	if len(os.Args) == 3 {
//		ipList := os.Args[1]
//		portList := os.Args[2]
//		ips, err := util.GetIpList(ipList)
//		ports, err := util.GetPorts(portList)
//		if err != nil {
//			fmt.Println(err)
//		}
//		for _, ip := range ips {
//			for _, port := range ports {
//				_, _, err := scanner.Connect(ip.String(), port)
//				if err != nil {
//					continue
//				}
//				fmt.Printf("ip: %v, port: %v is open\n", ip, port)
//			}
//		}
//	} else {
//		fmt.Printf("%v iplist port\n", os.Args[0])
//	}
//}

func main() {
	if len(os.Args) == 3 {
		util.CheckRoot()

		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		_ = err

		for _, ip := range ips {
			for _, port := range ports {
				ip1, port1, err1 := scanner.SynScan(ip.String(), port)
				if err1 == nil && port1 > 0 {
					fmt.Printf("%v:%v is open\n", ip1, port1)
				}
			}
		}
	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}