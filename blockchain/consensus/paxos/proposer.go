package paxos

import (
	"log"
	"time"
)

type proposer struct {
	id int
	lastSeq int
	value string
	valueN int
	acceptors map[int]promise
	nt network
}

func newProposer(id int, value string, nt network, acceptors ...int) *proposer {
	p := &proposer{id: id, nt: nt, lastSeq: 0, value: value, acceptors: make(map[int]promise)}
	for _, a := range acceptors {
		p.acceptors[a] = message{}
	}

	return p
}

func (p *proposer) run() {
	var ok bool
	var m message

	for !p.majorityReached() {
		if !ok {
			ms := p.prepare()
			for i := range ms {
				p.nt.send(ms[i])
			}
		}

		m, ok = p.nt.recv(time.Second)
		if !ok {
			continue
		}
		switch m.typ {
		case Promise:
			p.receivePromise(m)
		default:
			log.Panicf("proposer: %d unexpected message type: %v", p.id, m.typ)
		}
	}
	log.Printf("proposer: %d promise %d reached majority %d", p.id, p.n(), p.majority())

	log.Printf("proposer: %d starts to propose [%d: %s]", p.id, p.n(), p.value)
	ms := p.propose()
	for i := range ms {
		p.nt.send(ms[i])
	}
}

func (p *proposer) propose() []message {
	ms := make([]message, p.majority())

	i := 0
	for to, promise := range p.acceptors {
		if promise.number() == p.n() {
			ms[i] = message{from: p.id, to: to, typ: Propose, n: p.n(), value: p.value}
			i++
		}
		if i == p.majority() {
			break
		}
	}

	return ms
}

func (p *proposer) prepare() []message {
	p.lastSeq++

	ms := make([]message, p.majority())
	i := 0
	for to := range p.acceptors {
		ms[i] = message{from: p.id, to: to, typ: Prepare, n: p.n()}
		i++
		if i == p.majority() {
			break
		}
	}

	return ms
}

func (p *proposer) receivePromise(promise message) {
	prevPromise := p.acceptors[promise.from]

	if prevPromise.number() < promise.number() {
		log.Printf("proposer: %d received a new promise %+v", p.id, promise)
		p.acceptors[promise.from] = promise

		if promise.proposalNumber() > p.valueN {
			log.Printf("proposer: %d updated the value [%s] to %s", p.id, p.value, promise.proposalValue())
			p.valueN = promise.proposalNumber()
			p.value = promise.proposalValue()
		}
	}
}

func (p *proposer) majority() int {
	return len(p.acceptors)/2 + 1
}

func (p *proposer) majorityReached() bool {
	m := 0
	for _, promise := range p.acceptors {
		if promise.number() == p.n() {
			m++
		}
	}
	if m >= p.majority() {
		return true
	}

	return false
}

func (p *proposer) n() int {
	return p.lastSeq<<16 | p.id
}

