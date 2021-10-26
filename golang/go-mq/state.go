package go_mq

type State int

const (
	Unacked State = iota
	Acked
	Rejected
	Pushed
)
