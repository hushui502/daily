package routers

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"honeypot/manager/models"
	"strconv"
)

func ListRecord(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	pre := p - 1
	if pre <= 0 {
		pre = 1
	}
	next := p + 1
	if sess.Get("admin") != nil {
		records, pages, total, _ := models.ListRecordByPage(p)
		pList := 0
		if pages-p > 10 {
			pList = p + 10
		} else {
			pList = pages
		}

		pageList := make([]int, 0)
		if pages <= 10 {
			for i := 1; i <= pList; i++ {
				pageList = append(pageList, i)
			}
		} else {
			if p <= 10 {
				for i := 1; i <= pList; i++ {
					pageList = append(pageList, i)
				}
			} else {
				t := p + 5
				if t > pages {
					t = pages
				}
				for i := p - 5; i <= t; i++ {
					pageList = append(pageList, i)
				}
			}
		}

		ctx.Data["total"] = total
		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["records"] = records
		ctx.HTML(200, "record")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ListRecordBySite(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	site := ctx.Params(":site")
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	pre := p - 1
	if pre <= 0 {
		pre = 1
	}
	next := p + 1
	if sess.Get("admin") != nil {
		records, pages, total, _ := models.ListRecordBySite(site, p)
		pList := 0
		if pages-p > 10 {
			pList = p + 10
		} else {
			pList = pages
		}

		pageList := make([]int, 0)
		if pages <= 10 {
			for i := 1; i <= pList; i++ {
				pageList = append(pageList, i)
			}
		} else {
			if p <= 10 {
				for i := 1; i <= pList; i++ {
					pageList = append(pageList, i)
				}
			} else {
				t := p + 5
				if t > pages {
					t = pages
				}
				for i := p - 5; i <= t; i++ {
					pageList = append(pageList, i)
				}
			}
		}

		ctx.Data["total"] = total
		ctx.Data["site"] = site
		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["records"] = records
		ctx.HTML(200, "record_site")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func RecordDetail(ctx *macaron.Context, sess session.Store) {
	id := ctx.Params(":id")
	if sess.Get("admin") != nil {
		record, _ := models.RecordDetail(id)
		ctx.Data["record"] = record
		ctx.HTML(200, "record_detail")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
