package settings

import (
	"gopkg.in/ini.v1"
	"honeypot/server/log"
)

var (
	Cfg      *ini.File
	SECRET   string
	HttpPort int
)

func init() {
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)
	if err != nil {
		log.Logger.Panicln(err)
	}
	SECRET = Cfg.Section("").Key("SECRET").MustString("SECRET_KEY")
	HttpPort = Cfg.Section("").Key("HTTP_PORT").MustInt(8080)
}
