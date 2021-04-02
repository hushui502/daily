package models

import (
	"fmt"
	"go-sec/analysis/server/settings"
	"go-sec/analysis/server/util"
	"upper.io/db.v3"
	"upper.io/db.v3/mongo"
)

var (
	DbConfig   DbCONF
	DbSettings db.ConnectionURL
	Session    db.Database
)

type DbCONF struct {
	DbType string
	DbHost string
	DbPort int64
	DbUser string
	DbPass string
	DbName string
}

func init() {
	cfg := settings.Cfg
	sec := cfg.Section("database")
	DbConfig.DbType = sec.Key("DB_TYPE").MustString("mongodb")
	DbConfig.DbHost = sec.Key("DB_HOST").MustString("127.0.0.1")
	DbConfig.DbPort = sec.Key("DB_PORT").MustInt64(27017)
	DbConfig.DbUser = sec.Key("DB_USER").MustString("user")
	DbConfig.DbPass = sec.Key("DB_PASS").MustString("password")
	DbConfig.DbName = sec.Key("DB_NAME").MustString("proxy_honeypot")

	_ = NewDbEngine()
}

func NewDbEngine() error {
	switch DbConfig.DbType {
	case "mysql":
		util.Log.Info("will support mysql")
	case "mongodb":
		DbSettings = mongo.ConnectionURL{Host: fmt.Sprintf("%v:%v", DbConfig.DbHost, DbConfig.DbPort),
			User: DbConfig.DbUser, Password: DbConfig.DbPass, Database: DbConfig.DbName}
		Session, err := mongo.Open(DbSettings)
		util.Log.Warningf("settings: %v, session: %v, err: %v\n", DbSettings, Session, err)
		if err != nil {
			util.Log.Panicf("Connect Database failed, err: %v", err)
		}
		util.Log.Infof("DB Type: %v, Connect err status: %v", DbConfig.DbType, Session.Ping())
	}

	return nil
}
