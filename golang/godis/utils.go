package godis

import (
	"godis/datastruct/dict"
	"godis/datastruct/lock"
)

func makeTestDB() *DB {
	return &DB{
		data:   dict.MakeConcurrent(1),
		ttlMap: dict.MakeConcurrent(ttlDictSize),
		locker: lock.Make(lockerSize),
	}
}
