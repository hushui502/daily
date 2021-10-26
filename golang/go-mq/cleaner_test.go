package go_mq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleaner_Clean(t *testing.T) {
	flushConn, err := OpenConnection("cleaner-flush", "tcp", "localhost:6379", 1, nil)
	assert.NoError(t, err)
	assert.NoError(t, flushConn.stopHeartbeat())
	assert.NoError(t, flushConn.flushDb())

	// todo
}
