package web

import (
	"fmt"
	"gopkg.in/macaron.v1"
	"honeypot/server/log"
	"honeypot/server/modules/web/routers"
	"honeypot/server/settings"
	"net/http"
)

func Start() {
	m := macaron.Classic()
	m.Use(macaron.Renderer)

	m.Get("/", routers.Index)
	m.Post("/api/send", routers.RecvData)
	log.Logger.Infof("start web server at: %v", settings.HttpPort)
	log.Logger.Debug(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", settings.HttpPort), m))
}
