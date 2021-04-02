package web

import (
	"fmt"
	"net/http"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"

	"go-sec/analysis/server/util"
	"go-sec/analysis/server/web/routers"
)

func RunWeb(ctx *cli.Context) (err error) {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
	m.Use(cache.Cacher())

	m.Get("/", routers.Index)
	m.Get("/http/", routers.HttpReq)
	m.Get("/dns/", routers.Dns)

	m.Post("/api/packet/", routers.SendPacket)
	m.Post("/api/http/", routers.SendHTML)
	m.Post("/api/dns/", routers.SendDns)

	if ctx.IsSet("host") {
		HTTP_HOST = ctx.String("host")
	}

	if ctx.IsSet("port") {
		HTTP_PORT = ctx.Int("port")
	}

	util.Log.Infof("run server on %v", fmt.Sprintf("%v:%v", HTTP_HOST, HTTP_PORT))
	err = http.ListenAndServe(fmt.Sprintf("%v:%v", HTTP_HOST, HTTP_PORT), m)

	return err
}
