package routers

import (
	"go-sec/analysis/server/models"
	"gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) {
	info, _ := models.ListEvilInfo()
	ctx.Data["info"] = info
	ctx.HTML(200, "index")
}

func HttpReq(ctx *macaron.Context) {
	info, _ := models.ListEvilHttpReq()
	ctx.Data["info"] = info
	ctx.HTML(200, "http_req")
}

func Dns(ctx *macaron.Context) {
	info, _ := models.ListEvilDns()
	ctx.Data["info"] = info
	ctx.HTML(200, "dns")
}
