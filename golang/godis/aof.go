package godis

import (
	"godis/interface/redis"
	"godis/redis/reply"
	"os"
	"sync"
)

type DB struct {
	aofChan chan *reply.MultiBulkReply

	aofFile *os.File

	aofFilename string

	aofReWriteChan chan *reply.NullBulkReply

	pausingAof sync.RWMutex
}

type extra struct {
	// 此命令是否需要持久化，比如 get 就不需要
	toPersist bool
	// 若 specialAof == nil 则将命令原样持久化，否则持久化 specialAof 中的指令， 比如 expire 这种命令就不能原样持久化
	specialAof []*reply.MultiBulkReply
}

type CmdFunc func(db *DB, args [][]byte) (redis.Reply, *extra)
