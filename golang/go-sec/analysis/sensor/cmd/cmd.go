package cmd

import (
	"github.com/urfave/cli"
	"go-sec/analysis/sensor/sensor"
)

var Start = cli.Command{
	Name:        "start",
	Usage:       "startup traffic-analysis sensor",
	Description: "startup traffic-analysis sensor",
	Action:      sensor.Start,
	Flags: []cli.Flag{
		boolFlag("debug, d", "debug mode"),
		stringFlag("filter,f", "", "setting filters"),
		intFlag("length,l", 1024, "setting snapshot Length"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
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
