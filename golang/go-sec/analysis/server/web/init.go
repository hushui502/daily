package web

import "go-sec/analysis/server/settings"

var (
	HTTP_HOST string
	HTTP_PORT int
)

func init() {
	cfg := settings.Cfg
	HTTP_HOST = cfg.Section("").Key("HTTP_HOST").MustString("127.0.0.1")
	HTTP_PORT = cfg.Section("").Key("HTTP_PORT").MustInt(8080)
}
