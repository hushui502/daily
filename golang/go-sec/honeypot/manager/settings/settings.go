package settings

import (
	"gopkg.in/ini.v1"
	"honeypot/manager/logger"
	"honeypot/manager/vars"
)

var (
	Cfg *ini.File
)

func init() {
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)

	if err != nil {
		logger.Logger.Panicln(err)
	}

	vars.HttpHost = Cfg.Section("").Key("HTTP_HOST").MustString("127.0.0.1")
	vars.HttpPort = Cfg.Section("").Key("HTTP_PORT").MustInt(8000)
}
