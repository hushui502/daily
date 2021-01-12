package plugins

import (
	"database/sql"
	"fmt"
	_ "github.com/netxfly/mysql"
	"go-sec/scanner/password_crack/models"
)

func ScanMysql(service models.Service) (result models.ScanResult, err error) {
	result.Service = service

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", service.Username,
		service.Password, service.Ip, service.Port, "mysql")
	db, err := sql.Open("mysql", dataSourceName)
	defer func() {
		if db != nil {
			_ = db.Close()
		}
	}()

	if err != nil {
		return result, err
	}
	err = db.Ping()
	if err != nil {
		return result, err
	}

	result.Result = true

	return result, err
}
