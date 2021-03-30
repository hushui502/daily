package main

import (
	"github.com/urfave/cli"
	"honeypot/agent/cmd"
	"honeypot/agent/util"
	"honeypot/agent/vars"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	vars.CurrentDir = util.GetCurDir()
	vars.CaKey = filepath.Join(vars.CurrentDir, "./certs/ca.key")
	vars.CaCert = filepath.Join(vars.CurrentDir, "./certs/ca.cert")
}

func main() {
	app := cli.NewApp()
	app.Usage = "proxy-agent"
	app.Version = "0.1"
	app.Author = "fft"
	app.Email = "fft@xmail.com"
	app.Commands = []cli.Command{cmd.Serve}
	app.Flags = append(app.Flags, cmd.Serve.Flags...)
	_ = app.Run(os.Args)
}
