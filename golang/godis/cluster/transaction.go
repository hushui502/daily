package cluster

import (
	"fmt"
	"godis"
	"godis/interface/redis"
	"godis/lib/logger"
	"godis/lib/timewheel"
	"godis/redis/reply"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Transaction struct {
	id      string   // transaction id
	args    [][]byte // cmd args
	cluster *Cluster
	conn    redis.Connection

	keys       []string // related keys
	lockedKeys bool
	undoLog    map[string][][]byte // store data for undoLog

	status int8
	mu     *sync.Mutex
}

const (
	maxLockTime       = 3 * time.Second
	waitBeforeCleanTx = 2 * maxLockTime

	createdStatus    = 0
	preparedStatus   = 1
	committedStatus  = 2
	rolledBackStatus = 3
)

func genTaskKey(txID string) string {
	return "tx:" + txID
}

func NewTransaction(cluster *Cluster, c redis.Connection, id string, args [][]byte, keys []string) *Transaction {
	return &Transaction{
		id:      id,
		args:    args,
		cluster: cluster,
		conn:    c,
		keys:    keys,
		status:  createdStatus,
		mu:      new(sync.Mutex),
	}
}

func (tx *Transaction) lockKeys() {
	if !tx.lockedKeys {
		tx.cluster.db.Locks(tx.keys...)
		tx.lockedKeys = true
	}
}

func (tx *Transaction) unlockKeys() {
	if tx.lockedKeys {
		tx.cluster.db.UnLocks(tx.keys...)
		tx.lockedKeys = false
	}
}

func (tx *Transaction) prepare() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	// lock keys
	tx.lockKeys()

	// build undolog
	tx.undoLog = make(map[string][][]byte)
	for _, key := range tx.keys {
		entity, ok := tx.cluster.db.GetEntity(key)
		if ok {
			blob := godis.EntityToCmd(key, entity)
			tx.undoLog[key] = blob.Args
		} else {
			tx.undoLog[key] = nil
		}
	}

	tx.status = preparedStatus
	taskKey := genTaskKey(tx.id)
	timewheel.Delay(maxLockTime, taskKey, func() {
		// rollback transaction uncommitted until expire
		if tx.status == preparedStatus {
			logger.Info("abort transaction: ", tx.id)
			_ = tx.rollback()
		}
	})

	return nil
}

func (tx *Transaction) rollback() error {
	curStatus := tx.status
	tx.mu.Lock()
	defer tx.mu.Unlock()

	// fast path: ensure tx status not changed by other goroutine
	if tx.status != curStatus {
		return fmt.Errorf("tx %s status changed", tx.id)
	}
	// no need to rollback a rolled-back transaction
	if tx.status == rolledBackStatus {
		return nil
	}

	tx.lockKeys()
	for key, blob := range tx.undoLog {
		if len(blob) > 0 {
			tx.cluster.db.Remove(key)
			tx.cluster.db.Exec(nil, blob)
		} else {
			tx.cluster.db.Remove(key)
		}
	}
	tx.unlockKeys()
	tx.status = rolledBackStatus

	return nil
}

// Rollback rollbacks local transaction
func Rollback(cluster *Cluster, c redis.Connection, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rollback' command")
	}
	txID := string(args[1])
	raw, ok := cluster.transactions.Get(txID)
	if !ok {
		return reply.MakeIntReply(0)
	}
	tx, _ := raw.(*Transaction)
	err := tx.rollback()
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}
	// clean transaction
	timewheel.Delay(waitBeforeCleanTx, "", func() {
		cluster.transactions.Remove(tx.id)
	})
	return reply.MakeIntReply(1)
}

func commit(cluster *Cluster, c redis.Connection, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'commit' command")
	}
	txID := string(args[1])
	raw, ok := cluster.transactions.Get(txID)
	if !ok {
		return reply.MakeIntReply(0)
	}
	tx, _ := raw.(*Transaction)

	tx.mu.Lock()
	defer tx.mu.Unlock()

	cmd := strings.ToLower(string(tx.args[0]))
	var result redis.Reply
	if cmd == "del" {
		result = commitDel(cluster, c, tx)
	} else if cmd == "mset" {
		result = commitMSet(cluster, c, tx)
	}

	if reply.IsErrorReply(result) {
		err2 := tx.rollback()
		return reply.MakeErrReply(fmt.Sprintf("err occurs when rollback: %v, origin err: %s", err2, result))
	}

	// after committed
	tx.unlockKeys()
	tx.status = committedStatus
	// clean finished transaction
	// do not clean immediately, in case rollback
	timewheel.Delay(waitBeforeCleanTx, "", func() {
		cluster.transactions.Remove(tx.id)
	})

	return result
}

func requestCommit(cluster *Cluster, c redis.Connection, txID int64, peers map[string][]string) ([]redis.Reply, reply.ErrorReply) {
	var errReply reply.ErrorReply
	txIDStr := strconv.FormatInt(txID, 10)
	respList := make([]redis.Reply, 0, len(peers))
	for peer := range peers {
		var resp redis.Reply
		if peer == cluster.self {
			resp = commit(cluster, c, makeArgs("commit", txIDStr))
		} else {
			resp = cluster.relay(peer, c, makeArgs("commit", txIDStr))
		}
		if reply.IsErrorReply(resp) {
			errReply = resp.(reply.ErrorReply)
			break
		}
		respList = append(respList, resp)
	}
	if errReply != nil {
		requestRollback(cluster, c, txID, peers)
		return nil, errReply
	}
	return respList, nil
}

func requestRollback(cluster *Cluster, c redis.Connection, txID int64, peers map[string][]string) {
	txIDStr := strconv.FormatInt(txID, 10)
	for peer := range peers {
		if peer == cluster.self {
			Rollback(cluster, c, makeArgs("rollback", txIDStr))
		} else {
			cluster.relay(peer, c, makeArgs("rollback", txIDStr))
		}
	}
}
