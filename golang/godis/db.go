package godis

import (
	"fmt"
	"godis/config"
	"godis/datastruct/dict"
	"godis/datastruct/lock"
	"godis/interface/redis"
	"godis/lib/logger"
	"godis/lib/timewheel"
	"godis/pubsub"
	"godis/redis/reply"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type DataEntity struct {
	Data interface{}
}

const (
	dataDictSize = 1 << 16
	ttlDictSize  = 1 << 10
	lockerSize   = 128
	aofQueueSize = 1 << 16
)

type cmdFunc func(db *DB, args [][]byte) redis.Reply

type DB struct {
	// key -> DataEntry
	data dict.Dict
	// key -> expireTime(time.Time)
	ttlMap dict.Dict

	// use this mutex for complicated command only, eg. rpush, incr...
	locker *lock.Locks
	// stop all data access for FlushDB
	stopWorld sync.WaitGroup
	// handle publish/subscribe
	hub *pubsub.Hub

	// main goroutine send commands to aof goroutines through aofChan
	aofChan     chan *reply.MultiBulkReply
	aofFile     *os.File
	aofFilename string

	// aof goroutines will send msg to main goroutines through this channel
	// when aof takes finished and ready to shutdown
	aofFinished chan struct{}
	// buffer commands received during aof rewrite progress
	aofRewriteBuffer chan *reply.MultiBulkReply
	// pause aof for start/finish aof rewrite progress
	pausingAof sync.RWMutex
}

var router = makeRouter()

func MakeDB() *DB {
	db := &DB{
		data:   dict.MakeConcurrent(dataDictSize),
		ttlMap: dict.MakeConcurrent(dataDictSize),
		locker: lock.Make(lockerSize),
		hub:    pubsub.MakeHub(),
	}

	// aof
	if config.Properties.AppendOnly {
		db.aofFilename = config.Properties.AppendFilename
		db.loadAof(0)
		aofFile, err := os.OpenFile(db.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			logger.Warn(err)
		} else {
			db.aofFile = aofFile
			db.aofChan = make(chan *reply.MultiBulkReply, aofQueueSize)
		}
		db.aofFinished = make(chan struct{})
		go func() {
			db.handleAof()
		}()
	}

	return db
}

func (db *DB) Close() {
	if db.aofFile != nil {
		close(db.aofChan)
		<-db.aofFinished
		err := db.aofFile.Close()
		if err != nil {
			logger.Warn(err)
		}
	}
}

func (db *DB) Exec(c redis.Connection, cmdArgs [][]byte) (result redis.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occur: %v\n%s", err, string(debug.Stack())))
			result = &reply.UnknownErrReply{}
		}
	}()

	cmd := strings.ToLower(string(cmdArgs[0]))
	if cmd == "auth" {
		return Auth(db, c, cmdArgs[1:])
	}
	if !isAuthenticated(c) {
		return reply.MakeErrReply("NOAUTH Authentication required")
	}

	// special commands
	if cmd == "subscribe" {
		if len(cmdArgs) < 2 {
			return &reply.ArgNumErrReply{Cmd: "subscribe"}
		}
		return pubsub.Subscribe(db.hub, c, cmdArgs[1:])
	} else if cmd == "publish" {
		return pubsub.Publish(db.hub, cmdArgs[1:])
	} else if cmd == "unsubscribe" {
		return pubsub.UnSubscribe(db.hub, c, cmdArgs[1:])
	} else if cmd == "bgrewriteaof" {
		// aof.go imports router.go, router.go cannot import BGRewriteAOF from aof.go
		return BGRewriteAOF(db, cmdArgs[1:])
	}

	// normal commands
	fun, ok := router[cmd]
	if !ok {
		return reply.MakeErrReply("ERR unknown command '" + cmd + "'")
	}
	if len(cmdArgs) > 1 {
		result = fun(db, cmdArgs[1:])
	} else {
		result = fun(db, [][]byte{})
	}
	return
}

/*		data access  	*/

func (db *DB) GetEntity(key string) (*DataEntity, bool) {
	db.stopWorld.Wait()

	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	if db.IsExpired(key) {
		return nil, false
	}
	entity, _ := raw.(*DataEntity)

	return entity, true
}

func (db *DB) PutEntity(key string, entity *DataEntity) int {
	db.stopWorld.Wait()

	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *DataEntity) int {
	db.stopWorld.Wait()
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *DataEntity) int {
	db.stopWorld.Wait()
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(keys ...string) (deleted int) {
	db.stopWorld.Wait()
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.data.Remove(key)
			db.ttlMap.Remove(key)
			deleted++
		}
	}

	return deleted
}

// Removes the given keys from db
func (db *DB) Removes(keys ...string) (deleted int) {
	db.stopWorld.Wait()
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.data.Remove(key)
			db.ttlMap.Remove(key)
			deleted++
		}
	}
	return deleted
}

func (db *DB) Flush() {
	db.stopWorld.Add(1)
	defer db.stopWorld.Done()

	db.data = dict.MakeConcurrent(dataDictSize)
	db.ttlMap = dict.MakeConcurrent(ttlDictSize)
	db.locker = lock.Make(lockerSize)
}

/*		lock function 	*/
// Lock locks key for writing (exclusive lock)
func (db *DB) Lock(key string) {
	db.locker.Lock(key)
}

// RLock locks key for read (shared lock)
func (db *DB) RLock(key string) {
	db.locker.RLock(key)
}

// UnLock release exclusive lock
func (db *DB) UnLock(key string) {
	db.locker.UnLock(key)
}

// RUnLock release shared lock
func (db *DB) RUnLock(key string) {
	db.locker.RUnLock(key)
}

// Locks lock keys for writing (exclusive lock)
func (db *DB) Locks(keys ...string) {
	db.locker.Locks(keys...)
}

// RLocks lock keys for read (shared lock)
func (db *DB) RLocks(keys ...string) {
	db.locker.RLocks(keys...)
}

// UnLocks release exclusive locks
func (db *DB) UnLocks(keys ...string) {
	db.locker.UnLocks(keys...)
}

// RUnLocks release shared locks
func (db *DB) RUnLocks(keys ...string) {
	db.locker.RUnLocks(keys...)
}

/* ---- TTL Functions ---- */

func genExpireTask(key string) string {
	return "expire:" + key
}

// Expire sets TTL of key
func (db *DB) Expire(key string, expireTime time.Time) {
	db.stopWorld.Wait()
	db.ttlMap.Put(key, expireTime)
	taskKey := genExpireTask(key)
	timewheel.At(expireTime, taskKey, func() {
		logger.Info("expire " + key)
		db.Remove(key)
	})
}

func (db *DB) Persist(key string) {
	db.stopWorld.Wait()
	db.ttlMap.Remove(key)
	taskKey := genExpireTask(key)
	timewheel.Cancel(taskKey)
}

func (db *DB) IsExpired(key string) bool {
	rawExpireTime, ok := db.ttlMap.Get(key)
	if !ok {
		return false
	}
	expireTime, _ := rawExpireTime.(time.Time)
	expired := time.Now().After(expireTime)
	if expired {
		db.Remove(key)
	}

	return expired
}

func (db *DB) AfterClientClose(c redis.Connection) {
	pubsub.UnsubscribeAll(db.hub, c)
}
