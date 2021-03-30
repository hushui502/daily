package util

import (
	"fmt"
	"github.com/elazarl/goproxy"
	goproxy_html "github.com/elazarl/goproxy/ext/html"
	"github.com/urfave/cli"
	"honeypot/agent/log"
	"honeypot/agent/modules"
	"honeypot/agent/vars"
	"net/http"
)

func init() {

}

func Start(ctx *cli.Context) {
	if ctx.IsSet("debug") {
		vars.DebugMode = ctx.Bool("debug")
	}
	if ctx.IsSet("port") {
		vars.ProxyPort = ctx.Int("port")
	}

	err := SetCA()
	log.Logger.Infof("caKey: %v, caCert: %v, set ca err: %v", vars.CaKey, vars.CaCert, err)

	proxy := goproxy.NewProxyHttpServer()
	log.Logger.Infof("proxy Start success, Listening on %v:%v ", vars.ProxyHost, vars.ProxyPort)

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(modules.ReqHandlerFunc)
	proxy.OnResponse(goproxy_html.IsWebRelatedText).DoFunc(modules.RespHandlerFunc)

	proxy.Verbose = vars.DebugMode
	log.Logger.Info(http.ListenAndServe(fmt.Sprintf("%v:%v", vars.ProxyHost, vars.ProxyPort), proxy))
}
