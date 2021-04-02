package models

import "time"

type Dns struct {
	DnsType string `json:"dns_type"`
	DnsName string `json:"dns_name"`
	SrcIp   string `json:"src_ip"`
	DstIp   string `json:"dst_ip"`
}

type EvilDns struct {
	Id       int64
	Time     time.Time `bson:"time"`
	SensorIp string    `bson:"sensor_ip"`
	IsEvil   bool      `bson:"is_evil"`
	Data     Dns       `bson:"data"`
}

func NewEvilDns(sensorIp string, isEvil bool, dns Dns) (evilDns *EvilDns) {
	now := time.Now()
	return &EvilDns{SensorIp: sensorIp, Time: now, IsEvil: isEvil, Data: dns}
}

func (d *EvilDns) Insert() error {
	_, err := Session.Collection("dns").Insert(d)
	return err
}

func ListEvilDns() ([]EvilDns, error) {
	result := make([]EvilDns, 0)
	res := Session.Collection("dns").Find("-_id").OrderBy().Limit(500)
	err := res.All(&result)
	return result, err
}
