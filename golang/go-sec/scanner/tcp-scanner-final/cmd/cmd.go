package cmd

import (
	"github.com/urfave/cli"
	"go-sec/scanner/tcp-scanner-final/util"
)

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "start to scan port",
	Description: "start to scan port",
	Action:      util.Scan,
	Flags: []cli.Flag{
		stringFlag("iplist, i", "", "ip list"),
		stringFlag("port, p", "", "port list"),
		stringFlag("mode, m", "", "scan mode"),
		intFlag("timeout, t", 3, "timeout"),
		intFlag("concurrency, c", 1000, "concurrency"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

//func flag(tp interface{}, name, value, usage string) interface{} {
//	switch tp.(type) {
//	case cli.StringFlag:
//		return cli.StringFlag{
//
//		}
//	}
//	return nil
//}
