package main

import (
	"errors"
	"runtime"
)

// build info
var (
	Version     string
	BuildTime   string
	GoVersion   = runtime.Version()
	errCanceled = errors.New("action canceled")
	homedir     string
	pathConf    = ".tcli_config"
	pathHist    = ".tcli_history"
)

func main() {

}
