package godis

import (
	List "godis/datastruct/list"
	"godis/interface/redis"
	"godis/redis/reply"
	"strconv"
)

func (db *DB) getAsList(key string) (*List.LinkedList, reply.ErrorReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}
	bytes, ok := entity.Data.(*List.LinkedList)
	if !ok {
		return nil, &reply.WrongTypeErrReply{}
	}
	return bytes, nil
}

func (db *DB) getOrInitList(key string) (list *List.LinkedList, isNew bool, errReply reply.ErrorReply) {
	list, errReply = db.getAsList(key)
	if errReply != nil {
		return nil, false, errReply
	}
	isNew = false
	if list == nil {
		list = &List.LinkedList{}
		db.PutEntity(key, &DataEntity{
			Data: list,
		})
		isNew = true
	}
	return list, isNew, nil
}

// LIndex gets element of list at given list
func LIndex(db *DB, args [][]byte) redis.Reply {
	// parse args
	if len(args) != 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'lindex' command")
	}
	key := string(args[0])
	index64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	index := int(index64)

	db.RLock(key)
	defer db.RUnLock(key)

	// get entity
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return &reply.NullBulkReply{}
	}

	size := list.Len() // assert: size > 0
	if index < -1*size {
		return &reply.NullBulkReply{}
	} else if index < 0 {
		index = size + index
	} else if index >= size {
		return &reply.NullBulkReply{}
	}

	val, _ := list.Get(index).([]byte)
	return reply.MakeBulkReply(val)
}

// LLen gets length of list
func LLen(db *DB, args [][]byte) redis.Reply {
	// parse args
	if len(args) != 1 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'llen' command")
	}
	key := string(args[0])

	db.RLock(key)
	defer db.RUnLock(key)

	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeIntReply(0)
	}

	size := int64(list.Len())
	return reply.MakeIntReply(size)
}

// LPop removes the first element of list, and return it
func LPop(db *DB, args [][]byte) redis.Reply {
	// parse args
	if len(args) != 1 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'lindex' command")
	}
	key := string(args[0])

	// lock
	db.Lock(key)
	defer db.UnLock(key)

	// get data
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return &reply.NullBulkReply{}
	}

	val, _ := list.Remove(0).([]byte)
	if list.Len() == 0 {
		db.Remove(key)
	}
	db.AddAof(makeAofCmd("lpop", args))
	return reply.MakeBulkReply(val)
}

// LPush inserts element at head of list
func LPush(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'lpush' command")
	}
	key := string(args[0])
	values := args[1:]

	// lock
	db.Lock(key)
	defer db.UnLock(key)

	// get or init entity
	list, _, errReply := db.getOrInitList(key)
	if errReply != nil {
		return errReply
	}

	// insert
	for _, value := range values {
		list.Insert(0, value)
	}

	db.AddAof(makeAofCmd("lpush", args))
	return reply.MakeIntReply(int64(list.Len()))
}

// LPushX inserts element at head of list, only if list exists
func LPushX(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'lpushx' command")
	}
	key := string(args[0])
	values := args[1:]

	// lock
	db.Lock(key)
	defer db.UnLock(key)

	// get or init entity
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeIntReply(0)
	}

	// insert
	for _, value := range values {
		list.Insert(0, value)
	}
	db.AddAof(makeAofCmd("lpushx", args))
	return reply.MakeIntReply(int64(list.Len()))
}

// LRange gets elements of list in given range
func LRange(db *DB, args [][]byte) redis.Reply {
	// parse args
	if len(args) != 3 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'lrange' command")
	}
	key := string(args[0])
	start64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	start := int(start64)
	stop64, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	stop := int(stop64)

	// lock key
	db.RLock(key)
	defer db.RUnLock(key)

	// get data
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return &reply.EmptyMultiBulkReply{}
	}

	// compute index
	size := list.Len() // assert: size > 0
	if start < -1*size {
		start = 0
	} else if start < 0 {
		start = size + start
	} else if start >= size {
		return &reply.EmptyMultiBulkReply{}
	}
	if stop < -1*size {
		stop = 0
	} else if stop < 0 {
		stop = size + stop + 1
	} else if stop < size {
		stop = stop + 1
	} else {
		stop = size
	}
	if stop < start {
		stop = start
	}

	// assert: start in [0, size - 1], stop in [start, size]
	slice := list.Range(start, stop)
	result := make([][]byte, len(slice))
	for i, raw := range slice {
		bytes, _ := raw.([]byte)
		result[i] = bytes
	}
	return reply.MakeMultiBulkReply(result)
}

// LRem removes element of list at specified index
func LRem(db *DB, args [][]byte) redis.Reply {
	// parse args
	if len(args) != 3 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'lrem' command")
	}
	key := string(args[0])
	count64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	count := int(count64)
	value := args[2]

	// lock
	db.Lock(key)
	defer db.UnLock(key)

	// get data entity
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeIntReply(0)
	}

	var removed int
	if count == 0 {
		removed = list.RemoveAllByVal(value)
	} else if count > 0 {
		removed = list.RemoveByVal(value, count)
	} else {
		removed = list.ReverseRemoveByVal(value, -count)
	}

	if list.Len() == 0 {
		db.Remove(key)
	}
	if removed > 0 {
		db.AddAof(makeAofCmd("lrem", args))
	}

	return reply.MakeIntReply(int64(removed))
}

// LSet puts element at specified index of list
func LSet(db *DB, args [][]byte) redis.Reply {
	// parse args
	if len(args) != 3 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'lset' command")
	}
	key := string(args[0])
	index64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	index := int(index64)
	value := args[2]

	// lock
	db.Lock(key)
	defer db.UnLock(key)

	// get data
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeErrReply("ERR no such key")
	}

	size := list.Len() // assert: size > 0
	if index < -1*size {
		return reply.MakeErrReply("ERR index out of range")
	} else if index < 0 {
		index = size + index
	} else if index >= size {
		return reply.MakeErrReply("ERR index out of range")
	}

	list.Set(index, value)
	db.AddAof(makeAofCmd("lset", args))
	return &reply.OkReply{}
}

// RPop removes last element of list then return it
func RPop(db *DB, args [][]byte) redis.Reply {
	// parse args
	if len(args) != 1 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rpop' command")
	}
	key := string(args[0])

	// lock
	db.Lock(key)
	defer db.UnLock(key)

	// get data
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return &reply.NullBulkReply{}
	}

	val, _ := list.RemoveLast().([]byte)
	if list.Len() == 0 {
		db.Remove(key)
	}
	db.AddAof(makeAofCmd("rpop", args))
	return reply.MakeBulkReply(val)
}

// RPopLPush pops last element of list-A then insert it to the head of list-B
func RPopLPush(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rpoplpush' command")
	}
	sourceKey := string(args[0])
	destKey := string(args[1])

	// lock
	db.Locks(sourceKey, destKey)
	defer db.UnLocks(sourceKey, destKey)

	// get source entity
	sourceList, errReply := db.getAsList(sourceKey)
	if errReply != nil {
		return errReply
	}
	if sourceList == nil {
		return &reply.NullBulkReply{}
	}

	// get dest entity
	destList, _, errReply := db.getOrInitList(destKey)
	if errReply != nil {
		return errReply
	}

	// pop and push
	val, _ := sourceList.RemoveLast().([]byte)
	destList.Insert(0, val)

	if sourceList.Len() == 0 {
		db.Remove(sourceKey)
	}

	db.AddAof(makeAofCmd("rpoplpush", args))
	return reply.MakeBulkReply(val)
}

// RPush inserts element at last of list
func RPush(db *DB, args [][]byte) redis.Reply {
	// parse args
	if len(args) < 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rpush' command")
	}
	key := string(args[0])
	values := args[1:]

	// lock
	db.Lock(key)
	defer db.UnLock(key)

	// get or init entity
	list, _, errReply := db.getOrInitList(key)
	if errReply != nil {
		return errReply
	}

	// put list
	for _, value := range values {
		list.Add(value)
	}
	db.AddAof(makeAofCmd("rpush", args))
	return reply.MakeIntReply(int64(list.Len()))
}

// RPushX inserts element at last of list only if list exists
func RPushX(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rpush' command")
	}
	key := string(args[0])
	values := args[1:]

	// lock
	db.Lock(key)
	defer db.UnLock(key)

	// get or init entity
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeIntReply(0)
	}

	// put list
	for _, value := range values {
		list.Add(value)
	}
	db.AddAof(makeAofCmd("rpushx", args))

	return reply.MakeIntReply(int64(list.Len()))
}
