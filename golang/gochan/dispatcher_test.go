package gochan

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

type Manager struct {
	dispatcher *Dispatcher
}

func (m *Manager) Dispatch(objID int, task TaskFunc) error {
	return m.dispatcher.Dispatch(objID, task)
}

func (m *Manager) Close() {
	m.dispatcher.Close()
}

func TestDispatcher(t *testing.T) {
	gochanNum := 3
	bufferNum := 10
	manager := Manager{
		dispatcher: NewDispatcher(gochanNum, bufferNum),
	}

	for i := range make([]int, 6) {
		objID := i
		task1 := func() error {
			return errors.New("task " + strconv.Itoa(objID))
		}
		err := manager.Dispatch(objID, task1)
		assert.Nil(t, err)
	}
	for _ = range make([]int, 20) {
		objID := -1
		taskM3 := func() error {
			return errors.New("task " + strconv.Itoa(objID))
		}
		err := manager.Dispatch(objID, taskM3)
		assert.Nil(t, err)
	}

	time.Sleep(time.Second)

	manager.Close()
	err := manager.Dispatch(1, func() error { return nil })
	assert.NotNil(t, err)
}
