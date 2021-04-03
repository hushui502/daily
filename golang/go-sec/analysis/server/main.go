package main

import (
	"github.com/urfave/cli"
	"go-sec/analysis/server/cmd"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "traffic-analysis server"
	app.Author = "ff"
	app.Email = "ff@xsec.io"
	app.Version = "20201210"
	app.Usage = "traffic-analysis server"
	app.Commands = []cli.Command{cmd.Start}
	app.Flags = append(app.Flags, cmd.Start.Flags...)
	_ = app.Run(os.Args)
}
