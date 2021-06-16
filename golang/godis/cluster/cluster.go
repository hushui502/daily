package cluster

//
//import (
//	"context"
//	"fmt"
//	"github.com/jolestar/go-commons-pool/v2"
//	"godis"
//	"godis/config"
//	"godis/datastruct/dict"
//	"godis/interface/redis"
//	"godis/lib/consistenthash"
//	"godis/lib/idgenerator"
//	"godis/lib/logger"
//	"godis/redis/reply"
//	"runtime/debug"
//	"strings"
//)
//
//type Cluster struct {
//	self string
//
//	nodes          []string
//	peerPicker     *consistenthash.Map
//	peerConnection map[string]*pool.ObjectPool
//
//	db *godis.DB
//	// id -> transaction
//	transactions *dict.SimpleDict
//
//	idGenerator *idgenerator.IDGenerator
//}
//
//const (
//	replicas = 4
//	lockSize = 64
//)
//
//var allowFastTransaction = true
//
//func MakeCluster() *Cluster {
//	cluster := &Cluster{
//		self:           config.Properties.Self,
//		db:             godis.MakeDB(),
//		transactions:   dict.MakeSimple(),
//		peerPicker:     consistenthash.New(replicas, nil),
//		peerConnection: make(map[string]*pool.ObjectPool),
//
//		idGenerator: idgenerator.MakeGenerator(config.Properties.Self),
//	}
//	contains := make(map[string]struct{})
//	// self is a peer, so need to +1
//	nodes := make([]string, 0, len(config.Properties.Peers)+1)
//	for _, peer := range config.Properties.Peers {
//		if _, ok := contains[peer]; ok {
//			continue
//		}
//		contains[peer] = struct{}{}
//		nodes = append(nodes, peer)
//	}
//	nodes = append(nodes, config.Properties.Self)
//	cluster.peerPicker.AddNode(nodes...)
//	ctx := context.Background()
//	for _, peer := range config.Properties.Peers {
//		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
//			Peer: peer,
//		})
//	}
//	cluster.nodes = nodes
//
//	return cluster
//}
//
//type CmdFunc func(cluster *Cluster, c redis.Connection, cmdAndArgs [][]byte) redis.Reply
//
//func (cluster *Cluster) Close() {
//	cluster.db.Close()
//}
//
//var router = makeRouter()
//
//func isAuthenticated(c redis.Connection) bool {
//	if config.Properties.RequirePass == "" {
//		return true
//	}
//	return c.GetPassword() == config.Properties.RequirePass
//}
//
//func (cluster *Cluster) Exec(c redis.Connection, cmdArgs [][]byte) (result redis.Reply) {
//	defer func() {
//		if err := recover(); err != nil {
//			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
//			result = &reply.UnknownErrReply{}
//		}
//	}()
//	cmd := strings.ToLower(string(cmdArgs[0]))
//	if cmd == "auth" {
//		return godis.Auth(cluster.db, c, cmdArgs[1:])
//	}
//	if !isAuthenticated(c) {
//		return reply.MakeErrReply("NOAUTH Authentication required")
//	}
//	cmdFunc, ok := router[cmd]
//	if !ok {
//		return reply.MakeErrReply("ERR unknown command '" + cmd + "', or not supported in cluster mode")
//	}
//	result = cmdFunc(cluster, c, cmdArgs)
//	return
//}
//
//func (cluster *Cluster) AfterClientClose(c redis.Connection) {
//	cluster.db.AfterClientClose(c)
//}
//
//func ping(cluster *Cluster, c redis.Connection, args [][]byte) redis.Reply {
//	return godis.Ping(cluster.db, args[1:])
//}
//
///*----- utils -------*/
//
//func makeArgs(cmd string, args ...string) [][]byte {
//	result := make([][]byte, len(args)+1)
//	result[0] = []byte(cmd)
//	for i, arg := range args {
//		result[i+1] = []byte(arg)
//	}
//	return result
//}
//
//// return peer -> keys
//func (cluster *Cluster) groupBy(keys []string) map[string][]string {
//	result := make(map[string][]string)
//	for _, key := range keys {
//		peer := cluster.peerPicker.PickNode(key)
//		group, ok := result[peer]
//		if !ok {
//			group = make([]string, 0)
//		}
//		group = append(group, key)
//		result[peer] = group
//	}
//	return result
//}
