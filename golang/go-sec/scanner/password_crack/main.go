package main

import (
	"github.com/urfave/cli"
	"go-sec/scanner/password_crack/cmd"
	"os"
	"runtime"
)

func main() {
	app := cli.NewApp()
	app.Name = "password-crack"
	app.Author = "ff"
	app.Email = "ff@gamil.com"
	app.Version = "2022/3/2"
	app.Usage = "Weak password scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	_ = err
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
