package plugins

import (
	"fmt"
	"github.com/go-redis/redis"
	"go-sec/scanner/password_crack/models"
	"go-sec/scanner/password_crack/vars"
)

func ScanRedis(s models.Service) (result models.ScanResult, err error) {
	result.Service = s
	opt := redis.Options{
		Addr:        fmt.Sprintf("%v:%v", s.Ip, s.Port),
		Password:    s.Password,
		DB:          0,
		DialTimeout: vars.Timeout,
	}
	client := redis.NewClient(&opt)
	defer func() {
		if client != nil {
			_ = client.Close()
		}
	}()

	_, err = client.Ping().Result()
	if err != nil {
		return result, err
	}

	result.Result = true

	return result, err
}
