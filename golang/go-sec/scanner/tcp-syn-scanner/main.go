package main

import (
	"fmt"
	"go-sec/scanner/tcp-syn-scanner/scanner"
	"go-sec/scanner/tcp-syn-scanner/util"
	"os"
	"runtime"
)

func main() {
	if len(os.Args) == 3 {
		util.CheckRoot()

		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		if err != nil {
			fmt.Println(err)
		}
		tasks, _ := scanner.GenerateTask(ips, ports)
		scanner.RunTask(tasks)
		scanner.PrintResult()
	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
