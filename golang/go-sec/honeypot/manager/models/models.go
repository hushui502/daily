package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"honeypot/manager/logger"
	"honeypot/manager/settings"
	"time"
)

var (
	Session  *mgo.Session
	Host     string
	Port     int
	USERNAME string
	PASSWORD string
	DataName string

	collAdmin *mgo.Collection
)

func init() {
	cfg := settings.Cfg
	sec := cfg.Section("MONGODB")
	Host = sec.Key("HOST").MustString("127.0.0.1")
	Port = sec.Key("PORT").MustInt(27017)
	USERNAME = sec.Key("USER").MustString("xproxy")
	PASSWORD = sec.Key("PASS").MustString("passw0rd")
	DataName = sec.Key("DATA").MustString("xproxy")
	err := NewMongodbClient()
	err = Session.Ping()
	logger.Logger.Infof("CONNECT MONGODB, err: %v", err)

	collAdmin = Session.DB(DataName).C("users")
	userCount, _ := collAdmin.Find(nil).Count()
	if userCount == 0 {
		_ = NewUser("xproxy", "x@xsec.io")
	}
}

// return a mongodb session
func NewMongodbClient() (err error) {
	url := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", USERNAME, PASSWORD, Host, Port, DataName)
	Session, err = mgo.Dial(url)
	if err == nil {
		Session.SetSocketTimeout(1 * time.Hour)
	} else {
		logger.Logger.Panicf("connect mongodb failed, url: %v, err: %v", url, err)
	}
	return err
}
