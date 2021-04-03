package models

import (
	"net/http"
	"net/url"
	"time"
)

type HttpReq struct {
	Host          string
	Ip            string
	Client        string
	Port          string
	URL           *url.URL
	Header        http.Header
	RequestURI    string
	Method        string
	ReqParameters url.Values
}

type EvilHttpReq struct {
	Id       int64
	Time     time.Time `bson:"time"`
	SensorIp string    `bson:"sensor_ip"`
	IsEvil   bool      `bson:"is_evil"`
	Data     HttpReq   `bson:"data"`
}

func NewEvilHttpReq(sensorIp string, isEvil bool, req HttpReq) (evilHttpReq *EvilHttpReq) {
	now := time.Now()
	return &EvilHttpReq{SensorIp: sensorIp, Time: now, IsEvil: isEvil, Data: req}
}

func (e *EvilHttpReq) Insert() {
	Session.Collection("http_req").Insert(e)
}

func ListEvilHttpReq() ([]EvilHttpReq, error) {
	evilHttpReqs := make([]EvilHttpReq, 0)
	res := Session.Collection("http_req").Find("-_id").OrderBy().Limit(500)
	err := res.All(&evilHttpReqs)
	return evilHttpReqs, err
}
