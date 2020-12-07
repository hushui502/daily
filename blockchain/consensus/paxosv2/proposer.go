package paxosv2

import (
	"sync"
	"sync/atomic"
	"time"
)

type Proposer struct {
	mu sync.Mutex
	me int	// proposer编号，提案者1， 2， 3...
	acceptorPeers []string	// acceptor的addr
	proposerID float32 // 提案号
	currentValue interface{}
	highestAcceptedID float32 // 收到acceptor返回的最高的接受的proposerID
	highestAcceptedValue interface{}
	decidedValue interface{}
}

func (p *Proposer) close() {}

func (p *Proposer) clean() {
	p.proposerID = 0
	p.currentValue = nil
	p.highestAcceptedID = 0
	p.highestAcceptedValue = nil
	p.decidedValue = nil
}

func (p *Proposer) propose(value interface{}) interface{} {
	// 不需要加锁
	p.currentValue = value
	p.runTwoPhase()
	return nil
}

func (p *Proposer) getMe() int {
	p.mu.Lock()
	me := p.me
	p.mu.Unlock()

	return me
}

func (p *Proposer) getQuorumSize() int64 {
	p.mu.Lock()
	size := int64(len(p.acceptorPeers)/2 + 1)
	p.mu.Unlock()

	return size
}

func (p *Proposer) getAcceptorPeers() []string {
	p.mu.Lock()
	peers := p.acceptorPeers
	p.mu.Unlock()

	return peers
}

func (p *Proposer) runTwoPhase() {
	peers := p.getAcceptorPeers()

	for p.decidedValue == nil {
		// phase 1
		prepareMsgReq := p.prepare()
		var promiseSuccessNum int64
		var promiseFailedNum int64
		promiseSuccessChan := make(chan struct{})
		promiseFailedChan := make(chan struct{})

		logPrint("[propose:%d] phase-1 , prepareMsg:%v", p.getMe(), prepareMsgReq)
		for _, peerAddr := range peers {
			go func(peerAddr string, prepareMsgReq PrepareMsg) {
				defer func() {
					if atomic.LoadInt64(&promiseSuccessNum) >= p.getQuorumSize() {
						promiseSuccessChan <- struct{}{}
						return
					}
					if atomic.LoadInt64(&promiseFailedNum) >= p.getQuorumSize() {
						promiseFailedChan <- struct{}{}
						return
					}
				}()

				promiseMsgResp, err := p.sendPrepare(peerAddr, &prepareMsgReq)
				if err != nil || promiseMsgResp.ProposeID != p.proposerID ||
					promiseMsgResp.AcceptorAddr != peerAddr {
					atomic.AddInt64(&promiseFailedNum, 1)
					return
				}

				if promiseMsgResp.Success {
					atomic.AddInt64(&promiseSuccessNum, 1)
				} else {
					atomic.AddInt64(&promiseFailedNum, 1)
				}

				if promiseMsgResp.AcceptedID > 0 {
					p.setAccepted(promiseMsgResp)
				}
			}(peerAddr, prepareMsgReq)
		}

		select {
		case <-time.After(200 * time.Second):
			continue
		case <-promiseFailedChan:
			continue
		case <-promiseSuccessChan:
			// prepare success
		}

		// phase 2
		var acceptSuccessNum int64
		var acceptFailedNum int64
		acceptSuccessChan := make(chan struct{})
		acceptFailedChan := make(chan struct{})

		// 如果有最高的值就会替换当前的current值
		acceptMsgReq := p.accept()
		logPrint("[proposer:%d] phase-2, acceptMsg:%v", p.getMe(), acceptMsgReq)

		for _, peerAddr := range peers {
			go func(peerAddr string, acceptMsgReq AcceptMsg) {
				defer func() {
					if atomic.LoadInt64(&acceptSuccessNum) > p.getQuorumSize() {
						acceptSuccessChan <- struct{}{}
						return
					}
					if atomic.LoadInt64(&acceptFailedNum) > p.getQuorumSize() {
						acceptFailedChan <- struct{}{}
						return
					}
				}()

				acceptedMsgResp, err := p.sendAccept(peerAddr, &acceptMsgReq)
				if err != nil || acceptMsgReq.ProposeID != acceptedMsgResp.ProposeID ||
					peerAddr != acceptedMsgResp.AcceptorAddr {
					atomic.AddInt64(&acceptFailedNum, 1)
					return
				}
				if acceptedMsgResp.Success {
					atomic.AddInt64(&acceptSuccessNum, 1)
				} else {
					atomic.AddInt64(&acceptFailedNum, 1)
				}
			}(peerAddr, acceptMsgReq)
		}
		select {
		case <-time.After(200 * time.Millisecond):
			continue
		case <-acceptFailedChan:
			continue
		case <-acceptSuccessChan:
			p.decidedValue = acceptMsgReq.Value
			logPrint("[proposer:%d] reach consensuse Value: %v", p.getMe(), acceptMsgReq.Value)
			return
		}
	}
}

func (p *Proposer) prepare() PrepareMsg {
	p.mu.Lock()
	proposerID := generateNumber(p.me, p.proposerID)
	msg := PrepareMsg{ProposeID:proposerID}
	p.proposerID = proposerID
	p.mu.Unlock()

	return msg
}

func (p *Proposer) accept() AcceptMsg {
	p.mu.Lock()
	msg := AcceptMsg{
		ProposeID:    p.proposerID,
		Value:        p.currentValue,
	}
	if p.highestAcceptedValue != nil {
		msg.Value = p.highestAcceptedValue
	}
	p.mu.Unlock()

	return msg
}

func (p *Proposer) setAccepted(promiseMsgResp *PromiseMsg) {
	p.mu.Lock()
	if promiseMsgResp.AcceptedID > p.highestAcceptedID {
		p.highestAcceptedID = promiseMsgResp.AcceptedID
		p.highestAcceptedValue = promiseMsgResp.AcceptedValue
	}
	p.mu.Unlock()
}

func (p *Proposer) sendPrepare(peerAddr string, msg *PrepareMsg) (*PromiseMsg, error) {
	reply := new(PromiseMsg)
	err := callRpc(peerAddr, "Acceptor", "RecievePrepare", msg, reply)

	return reply, err
}

func (p *Proposer) sendAccept(peerAddr string, msg *AcceptMsg) (*AcceptedMsg, error) {
	reply := new(AcceptMsg)
	err := callRpc(peerAddr, "Acceptor", "RecieveAccept", msg, reply)

	return reply, err
}
