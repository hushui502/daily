package util

import (
	"bufio"
	"go-sec/scanner/password_crack/logger"
	"go-sec/scanner/password_crack/models"
	"go-sec/scanner/password_crack/vars"
	"os"
	"strconv"
	"strings"
)

func ReadIpList(filename string) (ipList []models.IpAddr) {
	ipListFile, err := os.Open(filename)
	if err != nil {
		logger.Log.Fatalf("Open ip list file err, %v", err)
	}

	defer func() {
		if ipListFile != nil {
			_ = ipListFile.Close()
		}
	}()

	scanner := bufio.NewScanner(ipListFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		ipPort := strings.TrimSpace(line)
		t := strings.Split(ipPort, ":")
		ip := t[0]
		portProtocol := t[1]
		tmpPort := strings.Split(portProtocol, "|")
		// ip列表中指定了端口对应的服务
		if len(tmpPort) == 2 {
			port, _ := strconv.Atoi(tmpPort[0])
			protocol := strings.ToUpper(tmpPort[1])
			if vars.SupportProtocols[protocol] {
				addr := models.IpAddr{Ip: ip, Port: port, Protocol: protocol}
				ipList = append(ipList, addr)
			} else {
				logger.Log.Infof("Not support %v, ignore: %v:%v\n", protocol, ip, port)
			}
		} else {
			port, err := strconv.Atoi(tmpPort[0])
			if err != nil {
				protocol, ok := vars.PortNames[port]
				if ok && vars.SupportProtocols[protocol] {
					addr := models.IpAddr{Ip: ip, Port: port, Protocol: protocol}
					ipList = append(ipList, addr)
				}
			}
		}
	}

	return ipList
}

func ReadUserDict(userDict string) (users []string, err error) {
	file, err := os.Open(userDict)
	if err != nil {
		logger.Log.Fatalf("Open user dict file err: %v\n", err)
	}

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}

	return users, err
}

func ReadPasswordDict(passDict string) (passwords []string, err error) {
	file, err := os.Open(passDict)
	if err != nil {
		logger.Log.Fatalf("Open password dic file err: %v\n", err)
	}

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		password := strings.TrimSpace(scanner.Text())
		if passDict != "" {
			passwords = append(passwords, password)
		}
	}

	return passwords, err
}
