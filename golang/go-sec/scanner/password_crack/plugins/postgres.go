package plugins

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-sec/scanner/password_crack/models"
)

func ScanPostgres(service models.Service) (result models.ScanResult, err error) {
	result.Service = service

	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", service.Username,
		service.Password, service.Ip, service.Port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return result, err
	}
	defer func() {
		if db != nil {
			_ = db.Close()
		}
	}()

	err = db.Ping()
	if err != nil {
		return result, err
	}

	return result, err
}
