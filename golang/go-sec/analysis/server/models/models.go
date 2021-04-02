package models

import (
	"strings"
	"time"
	"upper.io/db.v3"
)

type ConnectionInfo struct {
	Protocol string `json:"protocol"`
	SrcIp    string `json:"src_ip"`
	SrcPort  string `json:"src_port"`
	DstIp    string `json:"dst_ip"`
	DstPort  string `json:"dst_port"`
}

type Source struct {
	Desc   string `json:"desc"`
	Source string `json:"source"`
}

type EvilIps struct {
	Ips []string `json:"ips"`
	Src Source   `json:"src"`
}

type IpList struct {
	Id   int64
	Ip   string   `json:"ip"`
	Info []Source `json:"info"`
}

type IplistApi struct {
	Evil bool   `json:"evil"`
	Data IpList `json:"data"`
}

type EvilConnectInfo struct {
	Id       int64
	Time     time.Time `bson:"time"`
	SensorIp string    `bson:"sensor_ip"`
	Protocol string    `bson:"protocol"`
	SrcIp    string    `bson:"src_ip"`
	SrcPort  string    `bson:"src_port"`
	DstIp    string    `bson:"dst_ip" `
	DstPort  string    `bson:"dst_port" `
	IsEvil   bool      `bson:"is_evil" `
	Data     []Source  `bson:"data"`
}

func NewEvilConnectionInfo(sensorIp string, info ConnectionInfo, evilData IplistApi) (evilInfo *EvilConnectInfo) {
	now := time.Now()
	return &EvilConnectInfo{SensorIp: sensorIp, Time: now, Protocol: info.Protocol, SrcIp: info.SrcIp,
		SrcPort: info.SrcPort, DstIp: info.DstIp, DstPort: info.DstPort, IsEvil: evilData.Evil, Data: evilData.Data.Info}
}

func (i *EvilConnectInfo) Insert() (err error) {
	err, IsExist := i.Exist()
	if err == nil && !IsExist {
		Session.Collection("connection_info").Insert(i)
	}
	return err
}

func (i *EvilConnectInfo) Exist() (err error, isExist bool) {
	srcIp := strings.Split(i.SrcIp, ":")[0]
	Cond := db.Cond{"src_ip": srcIp, "dst_ip": i.DstIp}
	res := Session.Collection("connection_info").Find(Cond)
	evilInfos := make([]EvilConnectInfo, 0)
	err = res.All(&evilInfos)
	// util.Log.Errorln(i.SrcIp, i.DstIp, evilInfos, len(evilInfos), err)
	if len(evilInfos) > 0 {
		isExist = true
	}
	return err, isExist
}

func ListEvilInfo() ([]EvilConnectInfo, error) {
	evilInfos := make([]EvilConnectInfo, 0)
	res := Session.Collection("connection_info").Find("-_id").OrderBy().Limit(500)
	err := res.All(&evilInfos)
	return evilInfos, err
}
