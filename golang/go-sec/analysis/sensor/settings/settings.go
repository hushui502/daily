package settings

import (
	"go-sec/analysis/sensor/misc"
	"gopkg.in/ini.v1"
)

var (
	Cfg        *ini.File
	DeviceName string
	DebugMode  bool
	FilterRule string
	Ips        []string
)

func init() {
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)
	if err != nil {
		misc.Log.Panicln(err)
	}

	DeviceName = Cfg.Section("").Key("DEVICE_NAME").MustString("eth0")
	DebugMode = Cfg.Section("").Key("DEBUG_MODE").MustBool(false)
	FilterRule = Cfg.Section("").Key("FILTER_RULE").MustString("tcp or (udp and dst port 53)")

	Ips, _ = misc.GetIpList(DeviceName)
	misc.Log.Infof("Device name:[%v], ip addr:%v, Debug mode:[%v]", DeviceName, Ips, DebugMode)
}
