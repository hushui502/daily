package main

import (
	"github.com/urfave/cli"
	"go-sec/analysis/sensor/cmd"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "traffic-analysis sensor"
	app.Author = "abc"
	app.Email = "abc@xsec.io"
	app.Version = "20211020"
	app.Usage = "traffic-analysis sensor"
	app.Commands = []cli.Command{cmd.Start}
	app.Flags = append(app.Flags, cmd.Start.Flags...)
	_ = app.Run(os.Args)
}
