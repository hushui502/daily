package main

import (
	"github.com/urfave/cli"
	"go-sec/scanner/tcp-scanner-final/cmd"
	"os"
	"runtime"
)

func main() {
	app := cli.NewApp()
	app.Name = "port_scanner"
	app.Author = "ff"
	app.Email = "ff@gmail.io"
	app.Version = "2020/2/3"
	app.Usage = "tcp syn/connect port scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	_ = err
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
