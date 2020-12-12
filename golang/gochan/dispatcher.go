package gochan

import (
	"errors"
	"math/rand"
	"sync/atomic"
)

const (
	dispacherStatusOpen int32 = 0
	dispacherStatusClosed int32 = 1
)

type Dispatcher struct {
	status int32
	gcNum int
	gcs []*gochan
}

func NewDispatcher(gochanNum, bufferNum int) *Dispatcher {
	logger.Infof("%d gochans and %d bufferNum chan buffer", gochanNum, bufferNum)

	d := &Dispatcher{
		gcNum:  gochanNum,
		gcs:    make([]*gochan, gochanNum),
		status: dispacherStatusOpen,
	}

	for index := range d.gcs {
		gc := newGochan(bufferNum)
		gc.setUUID(index)
		d.gcs[index] = gc
		gc.run()
	}

	return d
}

func (d *Dispatcher) Dispatch(objID int, task TaskFunc) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("dispatcher closed")
		}
	}()

	if objID < 0 {
		objID = rand.Intn(d.gcNum)
	}

	if atomic.LoadInt32(&d.status) == dispacherStatusClosed {
		return errors.New("dispatcher closed")
	}

	index := objID % d.gcNum
	d.gcs[index].taskChan <- task

	return
}

func (d *Dispatcher) Close() {
	if atomic.LoadInt32(&d.status) == dispacherStatusClosed {
		return
	}

	atomic.StoreInt32(&d.status, dispacherStatusClosed)

	for _, gc := range d.gcs {
		gc.dieChan <- struct{}{}
	}
}