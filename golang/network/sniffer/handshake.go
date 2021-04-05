package main

import "errors"

const (
	ClientHello       condition = 0x1
	ServerHello       condition = 0x2
	ClientKeyExchange condition = 0x10
	NewSessionTicket  condition = 0x4
)

var states = []condition{
	ClientHello,
	ServerHello,
	ClientKeyExchange,
	NewSessionTicket,
}

var ErrHSInvalidLength = errors.New("Invalid handshake length")

type condition uint8

type stackFunc struct {
	fn func([]byte, condition) (bool, error)
	condition
}

type Handshake struct {
	route string
	stack Stack
	ready bool
	synCh chan string
}

func NewHandShake(route string, sycnCh chan string) *Handshake {
	hs := &Handshake{
		route: route,
		synCh: sycnCh,
	}

	for i := len(states) - 1; i > 0; i-- {
		hs.stack.Push(&stackFunc{fn: checker, condition: states[i]})
	}

	return hs
}

func (h *Handshake) Unmarshal(b []byte) error {
	if h.ready {
		return nil
	}

	stFunc, ok := h.stack.Pop()
	if ok {
		fn, cond := stFunc.(*stackFunc).fn, stFunc.(*stackFunc).condition
		if ok, err := fn(b, cond); !ok || err != nil {
			return err
		}
	}

	if h.stack.IsEmpty() {
		h.ready = true
		h.sync()
	}

	return nil
}

func checker(b []byte, c condition) (bool, error) {
	if len(b[5:]) < 6 {
		return false, ErrHSInvalidLength
	}
	return condition(b[5:][0]) == c, nil
}

func (h *Handshake) sync() {
	h.synCh <- h.route
}
