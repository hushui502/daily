package models

import (
	"gopkg.in/mgo.v2/bson"
	"honeypot/manager/vars"
	"net/http"
	"net/url"
	"time"
)

type Password struct {
	Id                bson.ObjectId     `bson:"_id"`
	ResponseBody      string            `bson:"response_body"`
	RequestBody       string            `bson:"request_body"`
	DateStart         time.Time         `bson:"date_start"`
	URL               string            `bson:"url"`
	RequestParameters url.Values        `bson:"request_parameters"`
	FromIp            string            `bson:"from_ip"`
	Site              string            `bson:"site"`
	ResponseHeader    http.Header       `bson:"response_header"`
	RequestHeader     http.Header       `bson:"request_header"`
	Data              map[string]string `bson:"data"`
}

func ListPasswordByPage(page int) (passwords []Password, pages int, total int, err error) {
	coll := Session.DB(DataName).C("password")
	total, _ = coll.Find(nil).Count()

	if total%vars.PageSize == 0 {
		pages = total / vars.PageSize
	} else {
		pages = total/vars.PageSize + 1
	}

	if page >= pages {
		page = pages
	}
	if page < 1 {
		page = pages
	}

	nums := (page - 1) * vars.PageSize
	if nums < 0 {
		nums = 0
	}

	err = coll.Find(nil).Skip(nums).Limit(vars.PageSize).All(&passwords)
	return passwords, pages, total, err
}

func ListPasswordBySite(site string, page int) (passwords []Password, pages int, total int, err error) {
	coll := Session.DB(DataName).C("password")
	total, _ = coll.Find(bson.M{"site": site}).Count()

	if total%vars.PageSize == 0 {
		pages = total / vars.PageSize
	} else {
		pages = total/vars.PageSize + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}

	i := (page - 1) * vars.PageSize
	if i < 0 {
		i = 0
	}

	err = coll.Find(bson.M{"site": site}).Skip(i).Limit(vars.PageSize).All(&passwords)
	return passwords, pages, total, err
}

func PasswordDetail(id string) (Password, error) {
	var password Password
	coll := Session.DB(DataName).C("password")
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&password)
	return password, err
}
