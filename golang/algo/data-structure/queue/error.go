package queue

import "errors"

var (
	//
	ErrDisposed = errors.New("queue: disposed")
	ErrTimeout = errors.New("queue: poll timed out")
	ErrEmptyQueue = errors.New("queue: empty queue")
)


