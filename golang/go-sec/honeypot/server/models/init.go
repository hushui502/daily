package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"honeypot/server/log"
	"honeypot/server/settings"
	"upper.io/db.v3"
	"upper.io/db.v3/mongo"
)

var (
	DbConfig   DbCONF
	DbSettings db.ConnectionURL
	Session    db.Database

	Engine *xorm.Engine
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
	DbConfig.DbType = sec.Key("DB_TYPE").MustString("mysql")
	DbConfig.DbHost = sec.Key("DB_HOST").MustString("127.0.0.1")
	DbConfig.DbPort = sec.Key("DB_PORT").MustInt64(3306)
	DbConfig.DbUser = sec.Key("DB_USER").MustString("x-proxy")
	DbConfig.DbPass = sec.Key("DB_PASS").MustString("x@xsec.io")
	DbConfig.DbName = sec.Key("DB_NAME").MustString("x-proxy")

	_ = NewDbEngine()
}

func NewDbEngine() (err error) {
	switch DbConfig.DbType {
	case "mysql":
		dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
			DbConfig.DbUser, DbConfig.DbPass, DbConfig.DbHost, DbConfig.DbPort, DbConfig.DbName)
		engine, err := xorm.NewEngine("mysql", dataSourceName)
		if err == nil {
			err = engine.Ping()
			if err == nil {
				_ = engine.Sync2(new(Record))
			}
		}
	case "mongodb":
		_, _ = GetSession()
	}

	return err
}

func GetSession() (db.Database, error) {
	var err error
	if Session == nil {
		DbSettings = mongo.ConnectionURL{Host: fmt.Sprintf("%v:%v", DbConfig.DbHost, DbConfig.DbPort), User: DbConfig.DbUser,
			Password: DbConfig.DbPass, Database: DbConfig.DbName}
		Session, err = mongo.Open(DbSettings)
		if err != nil {
			log.Logger.Infof("DB Type: %v, DbSettings: %v, Connect err status: %v", DbConfig.DbType, DbSettings, Session.Ping())
		}
	}

	return Session, err
}
