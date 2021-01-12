package main

import (
	"github.com/urfave/cli"
	"go-sec/scanner/proxy-scanner/cmd"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "proxy scanner"
	app.Author = "ff"
	app.Email = "ff@gamil.com"
	app.Version = "2022/1/15"
	app.Usage = "A SOCKS4/SOCKS4a/SOCKS5/HTTP/HTTPS proxy scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	_ = app.Run(os.Args)
}
