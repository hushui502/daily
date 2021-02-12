package paxoskv

import (
	"context"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
	"time"
)

var (
	NotEnoughQuorum = errors.New("not enough quorum")
	AcceptorBasePort = int64(3333)
)

func (p *Proposer) RunPaxos(acceptorIds []int64, val *Value) *Value {
	quorum := len(acceptorIds)/2 + 1

	for {
		p.Val = nil
		maxVotedVal, higherBal, err := p.Phase1(acceptorIds, quorum)
		if err != nil {
			pretty.Logf("Proposer: fail to run phase-1: highest ballot: %v, increment ballot and retry", higherBal)
			p.Bal.N = higherBal.N + 1
			continue
		}

		if maxVotedVal == nil {
			pretty.Logf("Proposer: no voted value seen, propose my value: %v", val)
		} else {
			val = maxVotedVal
		}

		if val == nil {
			pretty.Logf("Proposer: no value to propose in phase-2, quit")
			return nil
		}

		p.Val = val
		pretty.Logf("Proposer: proposer chose value to propose: %s", p.Val)

		higerBal, err := p.Phase2(acceptorIds, quorum)
		if err != nil {
			p.Bal.N = higerBal.N + 1
			continue
		}

		pretty.Logf("Proposer: value is voted by a quorum and has been safe: %v", maxVotedVal)
		return p.Val
	}
}

type Version struct {
	mu sync.Mutex
	acceptor Acceptor
}

type Versions map[int64]*Version

type KVServer struct {
	mu sync.Mutex
	Storage map[string]Versions
}

func (s *KVServer) getLockedVersion(id *PaxosInstanceId) *Version {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := id.Key
	ver := id.Ver
	rec, found := s.Storage[key]
	if !found {
		rec = Versions{}
		s.Storage[key] = rec
	}

	v, found := rec[ver]
	if !found {
		// initialize an empty paxos instance
		rec[ver] = &Version{
			acceptor: Acceptor{
				LastBal: &BallotNum{},
				VBal:    &BallotNum{},
			},
		}
		v = rec[ver]
	}

	pretty.Logf("Acceptor: getLockedVersion: %s", v)
	v.mu.Lock()

	return v
}

func (s *KVServer) Prepare(c context.Context, r *Proposer) (*Acceptor, error) {
	pretty.Logf("Acceptor: recv Accept-request: %v", r)

	v := s.getLockedVersion(r.Id)
	defer v.mu.Unlock()

	reply := v.acceptor

	if r.Bal.GE(v.acceptor.LastBal) {
		v.acceptor.LastBal = r.Bal
	}

	return &reply, nil
}

func (s *KVServer) Accept(c context.Context, r *Proposer) (*Acceptor, error) {
	pretty.Logf("Acceptor: recv Accept-request: %v", r)

	v := s.getLockedVersion(r.Id)
	defer v.mu.Unlock()

	d := *v.acceptor.LastBal
	reply := Acceptor{
		LastBal: &d,
	}

	if r.Bal.GE(v.acceptor.LastBal) {
		v.acceptor.LastBal = r.Bal
		v.acceptor.Val = r.Val
		v.acceptor.VBal = r.Bal
	}

	return &reply, nil
}

// starts a grpc server for every acceptor.
func ServeAcceptors(acceptorIds []int64) []*grpc.Server {
	servers := []*grpc.Server{}

	for _, aid := range acceptorIds {
		addr := fmt.Sprintf(":%d", AcceptorBasePort+int64(aid))

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal("listen: %s %v", addr, err)
		}

		s := grpc.NewServer()
		RegisterPaxosKVServer(s, &KVServer{
			Storage: map[string]Versions{},
		})
		reflection.Register(s)
		pretty.Logf("Acceptor-%d serving on %s ...", aid, addr)
		servers = append(servers, s)
		go s.Serve(lis)
	}

	return servers
}


// GE compare two ballot number a, b and return whether a >= b
func (a *BallotNum) GE(b *BallotNum) bool {
	if a.N > b.N {
		return true
	}
	// 直观就好，交给编译器优化
	if a.N < b.N {
		return false
	}

	return a.ProposerId >= b.ProposerId
}

func (p *Proposer) Phase1(acceptorIds []int64, quorum int) (*Value, *BallotNum, error) {
	replies := p.rpcToAll(acceptorIds, "Prepare")

	ok := 0
	higherBal := *p.Bal
	maxVoted := &Acceptor{VBal: &BallotNum{}}

	for _, r := range replies {
		pretty.Logf("Proposer: handling Prepare reply: %s", r)
		if !p.Bal.GE(r.LastBal) {
			higherBal = *r.LastBal
			continue
		}
		// find the voted value with highest vbal
		if r.VBal.GE(maxVoted.VBal) {
			maxVoted = r
		}

		ok += 1
		if ok == quorum {
			return maxVoted.Val, nil, nil
		}
	}

	return nil, &higherBal, NotEnoughQuorum
}

func (p *Proposer) Phase2(acceptorIds []int64, quorum int) (*BallotNum, error) {
	replies := p.rpcToAll(acceptorIds, "Accept")

	ok := 0
	higherBal := *p.Bal
	for _, r := range replies {
		pretty.Logf("Proposer: handing Accept reply: %s", r)
		if !p.Bal.GE(r.LastBal) {
			higherBal = *r.LastBal
			continue
		}
		ok += 1
		if ok == quorum {
			return nil, nil
		}
	}

	return &higherBal, NotEnoughQuorum
}

// rpcToAll send prepare or accept rpc to the specified acceptors
func (p *Proposer) rpcToAll(acceptorIds []int64, action string) []*Acceptor {
	replies := []*Acceptor{}

	for _, aid := range acceptorIds {
		var err error
		address := fmt.Sprintf("127.0.0.1:%d", AcceptorBasePort+int64(aid))
		// set up a connection to the server
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		defer conn.Close()
		c := NewPaxosKVClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var reply *Acceptor
		if action == "Prepare" {
			reply, err = c.Prepare(ctx, p)
		} else if action == "Accept" {
			reply, err = c.Accept(ctx, p)
		}

		if err != nil {
			log.Printf("Proposer: %s failure from Acceptor-%d: %v", action, aid, err)
		}
		log.Printf("Proposer: recv %s reply from: Acceptor-%d: %v", action, aid, reply)

		replies = append(replies, reply)
	}

	return replies
}