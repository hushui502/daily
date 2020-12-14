package cache

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisConnCachePool struct {
	pool *redis.Pool
}

func NewRedisCachePool(address, password string, idleTimeout, cap, maxIdle int) *RedisConnCachePool {
	redisPool := &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := conn.Do("AUTH", password); err != nil {
					_ = conn.Close()
					return nil, err
				}
			}
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				panic(err)
			}

			return err
		},
		MaxIdle:         maxIdle,
		MaxActive:       cap,
		IdleTimeout:     time.Duration(idleTimeout)*time.Second,
		Wait:            true,
	}

	return &RedisConnCachePool{pool:redisPool}
}

