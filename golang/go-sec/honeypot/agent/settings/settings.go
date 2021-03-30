package settings

import (
	"gopkg.in/ini.v1"
	"log"
)

var (
	Cfg *ini.File
)

func init() {
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)

	if err != nil {
		log.Panicf("load conf/app.ini failed, err: %v\n", err)
	}
}
