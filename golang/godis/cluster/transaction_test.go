package cluster

import (
	"godis/redis/reply/asserts"
	"math/rand"
	"strconv"
	"testing"
)

func TestRollback(t *testing.T) {
	// rollback uncommitted transaction
	FlushAll(testCluster, nil, toArgs("FLUSHALL"))
	txID := rand.Int63()
	txIDStr := strconv.FormatInt(txID, 10)
	keys := []string{"a", "b"}
	groupMap := testCluster.groupBy(keys)
	args := []string{txIDStr}
	args = append(args, keys...)
	testCluster.Exec(nil, toArgs("SET", "a", "a"))
	ret := prepareDel(testCluster, nil, makeArgs("PrepareDel", args...))
	asserts.AssertNotError(t, ret)
	requestRollback(testCluster, nil, txID, groupMap)
	ret = testCluster.Exec(nil, toArgs("GET", "a"))
	asserts.AssertBulkReply(t, ret, "a")

	// rollback committed transaction
	FlushAll(testCluster, nil, toArgs("FLUSHALL"))
	txID = rand.Int63()
	txIDStr = strconv.FormatInt(txID, 10)
	args = []string{txIDStr}
	args = append(args, keys...)
	testCluster.Exec(nil, toArgs("SET", "a", "a"))
	ret = prepareDel(testCluster, nil, makeArgs("PrepareDel", args...))
	asserts.AssertNotError(t, ret)
	_, err := requestCommit(testCluster, nil, txID, groupMap)
	if err != nil {
		t.Errorf("del failed %v", err)
		return
	}
	ret = testCluster.Exec(nil, toArgs("GET", "a"))
	asserts.AssertNullBulk(t, ret)
	requestRollback(testCluster, nil, txID, groupMap)
	ret = testCluster.Exec(nil, toArgs("GET", "a"))
	asserts.AssertBulkReply(t, ret, "a")
}
