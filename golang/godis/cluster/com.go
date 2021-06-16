package cluster

import (
	"context"
	"errors"
	"godis/interface/redis"
	"godis/redis/client"
	"godis/redis/reply"
)

func (cluster *Cluster) getPeerClient(peer string) (*client.Client, error) {
	connectionFactory, ok := cluster.peerConnection[peer]
	if !ok {
		return nil, errors.New("connection factory not found")
	}
	raw, err := connectionFactory.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	conn, ok := raw.(*client.Client)
	if !ok {
		return nil, errors.New("connection factory make wrong type")
	}
	return conn, nil
}

func (cluster *Cluster) returnPeerClient(peer string, peerClient *client.Client) error {
	connectionFactory, ok := cluster.peerConnection[peer]
	if !ok {
		return errors.New("connection factory not found")
	}
	return connectionFactory.ReturnObject(context.Background(), peerClient)
}

func (cluster *Cluster) relay(peer string, c redis.Connection, args [][]byte) redis.Reply {
	if peer == cluster.self {
		return cluster.db.Exec(c, args)
	}

	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}
	defer func() {
		_ = cluster.returnPeerClient(peer, peerClient)
	}()

	return peerClient.Send(args)
}

func (cluster *Cluster) broadcast(c redis.Connection, args [][]byte) map[string]redis.Reply {
	result := make(map[string]redis.Reply)
	for _, node := range cluster.nodes {
		reply := cluster.relay(node, c, args)
		result[node] = reply
	}
	return result
}
