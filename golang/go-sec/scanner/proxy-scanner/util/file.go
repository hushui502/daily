package util

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type ProxyAddr struct {
	Ip string
	Port int
}

func ReadProxyAddr(fileName string) (proxyAddrSlice []ProxyAddr) {
	proxyFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Open proxy file err, %v\n", err)
	}
	defer proxyFile.Close()

	scanner := bufio.NewScanner(proxyFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ipPort := strings.TrimSpace(scanner.Text())

		if ipPort == "" {
			continue
		}

		t := strings.Split(ipPort, ":")
		ip := t[0]
		port, err := strconv.Atoi(t[1])
		if err != nil {
			log.Fatalf("Convert port failed(string->int) %s", err)
		}
		proxyAddr := ProxyAddr{Ip: ip, Port: port}
		proxyAddrSlice = append(proxyAddrSlice, proxyAddr)
	}

	return proxyAddrSlice
}
