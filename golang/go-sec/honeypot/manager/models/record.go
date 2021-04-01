package models

import (
	"gopkg.in/mgo.v2/bson"
	"honeypot/manager/vars"
	"net/http"
	"net/url"
	"time"
)

type (
	HttpRecord struct {
		Id            bson.ObjectId `bson:"_id"`
		Session       int64         `json:"session"`
		Method        string        `json:"method"`
		RemoteAddr    string        `json:"remote_addr" bson:"remote"`
		StatusCode    int           `json:"status"`
		ContentLength int64         `json:"content_length"`
		Host          string        `json:"host"`
		Port          string        `json:"port"`
		Url           string        `json:"url"`
		Scheme        string        `json:"scheme"`
		Path          string        `json:"path"`
		ReqHeader     http.Header   `json:"req_header"`
		RespHeader    http.Header   `json:"resp_header"`
		RequestParam  url.Values    `json:"request_param" bson:"requestparameters"`
		RequestBody   []byte        `json:"request_body"`
		ResponseBody  []byte        `json:"response_body"`
		VisitTime     time.Time     `json:"visit_time"`
	}
)

func ListRecordByPage(page int) (records []HttpRecord, pages int, total int, err error) {

	coll := Session.DB(DataName).C("record")
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
		page = 1
	}

	i := (page - 1) * vars.PageSize
	if i < 0 {
		i = 0
	}

	err = coll.Find(nil).Skip(i).Limit(vars.PageSize).All(&records)
	return records, pages, total, err
}

func ListRecordBySite(site string, page int) (records []HttpRecord, pages int, total int, err error) {

	coll := Session.DB(DataName).C("record")
	total, _ = coll.Find(bson.M{"host": site}).Count()

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

	err = coll.Find(bson.M{"host": site}).Skip(i).Limit(vars.PageSize).All(&records)
	return records, pages, total, err
}

func RecordDetail(id string) (HttpRecord, error) {
	var record HttpRecord
	coll := Session.DB(DataName).C("record")
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&record)
	return record, err
}
