package db

import "godis/interface/redis"

type DB interface {
	Exec(client redis.Connection, args [][]byte) redis.Reply
	AfterClientClose(r redis.Connection)
	Close()
}
