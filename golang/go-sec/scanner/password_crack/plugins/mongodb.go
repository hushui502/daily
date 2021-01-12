package plugins

import (
	"fmt"
	"go-sec/scanner/password_crack/models"
	"go-sec/scanner/password_crack/vars"
	"gopkg.in/mgo.v2"
)

func ScanMongodb(s models.Service) (result models.ScanResult, err error) {
	result.Service = s
	url := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", s.Username, s.Password, s.Ip, s.Port, "test")
	session, err := mgo.DialWithTimeout(url, vars.Timeout)
	if err != nil {
		defer session.Close()
		err = session.Ping()
		if err == nil {
			result.Result = true
		}
	}

	return result, err
}
