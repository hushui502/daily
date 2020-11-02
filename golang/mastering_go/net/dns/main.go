package main

import (
	"fmt"
	"net"
)

func lookIP(address string) ([]string, error) {
	hosts, err := net.LookupAddr(address)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func lookHostname(hostname string) ([]string, error) {
	IPs, _ := net.LookupHost(hostname)
	return IPs, nil
}


// go run hhttp.go 127.0.0.1/www.baidu.com
func main() {
	arguments := os.Args

	input := arguments[1]
	IPaddress := net.ParseIP(input)
	if IPaddress == nil {
		IPs, _ := lookHostname(input)
		for _, singleIP := range IPs {
			fmt.Println(singleIP)
		}
	} else {
		hosts, _ := lookIP(input)
		for _, hostName := range hosts {
			fmt.Println(hostName)
		}
	}

	getNS("mtsoukalos.eu")
}

func getNS(domain string) {
	NSs, _ := net.LookupNS(domain)
	for _, ns := range NSs {
		fmt.Println(ns.Host)
	}
}
