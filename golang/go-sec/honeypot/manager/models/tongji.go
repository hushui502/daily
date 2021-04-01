package models

import (
	"gopkg.in/mgo.v2/bson"
	"honeypot/manager/vars"
)

func TongjiPasswordBySite(page int) (passwords []bson.M, pages int, total int, err error) {
	coll := Session.DB(DataName).C("password")
	pipe := coll.Pipe([]bson.M{{"$group": bson.M{"_id": "$site", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}}})

	resp := []bson.M{}
	err = pipe.All(&resp)
	total = len(resp)

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

	if page*vars.PageSize-i > len(resp) {
		passwords = resp
	} else {
		passwords = resp[i : page*vars.PageSize]
	}
	return passwords, pages, total, err
}

func TongjiUrls(page int) (urls []bson.M, pages int, total int, err error) {
	coll := Session.DB(DataName).C("proxy_honeypot")
	pipe := coll.Pipe([]bson.M{{"$group": bson.M{"_id": "$host", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}}, {"allowDiskUse": true}})

	resp := []bson.M{}
	err = pipe.All(&resp)
	total = len(resp)

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

	urls = resp[i : page*vars.PageSize]
	return urls, pages, total, err
}
