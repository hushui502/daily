package paxosv2

import (
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
)

type Acceptor struct {
	mu            sync.Mutex
	localAddr     string
	learnerPeers  []string
	promiseID     float32
	acceptedID    float32
	acceptedValue interface{}
	listener      net.Listener
	isunreliable  bool
}

func (a *Acceptor) getLearnerPeers() []string {
	a.mu.Lock()
	peers := a.learnerPeers
	a.mu.Unlock()

	return peers
}

func (a *Acceptor) getAddr() string {
	a.mu.Lock()
	addr := a.localAddr
	a.mu.Unlock()

	return addr
}

func (a *Acceptor) RecievePrepare(arg *PrepareMsg, reply *PromiseMsg) error {
	logPrint("[acceptor %s RecievePrepare:%v ]", a.localAddr, arg)
	reply.ProposeID = arg.ProposeID
	reply.AcceptorAddr = a.getAddr()
	if arg.ProposeID > a.promiseID {
		a.promiseID = arg.ProposeID
		reply.Success = true
		if a.acceptedID > 0 && a.acceptedValue != nil {
			reply.AcceptedID = a.acceptedID
			reply.AcceptedValue = a.acceptedValue
		}
	}

	return nil
}

func (a *Acceptor) RecieveAccept(arg *AcceptMsg, reply *AcceptedMsg) error {
	reply.ProposeID = arg.ProposeID
	if arg.ProposeID == a.promiseID {
		reply.Success = true
		reply.AcceptorAddr = a.getAddr()
		a.promiseID = arg.ProposeID
		a.acceptedID = arg.ProposeID
		a.acceptedValue = arg.Value
		for _, learnerPeer := range a.getLearnerPeers() {
			callRpc(learnerPeer, "Learner", "RecieveAccepted", reply, &EmptyMsg{})
		}
	}

	return nil
}

func (a *Acceptor) startRpc() {
	rpcx := rpc.NewServer()
	rpcx.Register(a)
	l, err := net.Listen("tcp", a.localAddr)
	a.listener = l
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	//defer l.Close()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			if a.isunreliable && rand.Int63()%100 < 300 {
				conn.Close()
				continue
			}
			go rpcx.ServeConn(conn)
		}
	}()
}

func (a *Acceptor) clean() {
	a.promiseID = 0
	a.acceptedID = 0
	a.acceptedValue = nil
}

func (a *Acceptor) close() {
	a.listener.Close()
}
