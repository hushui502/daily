package gochan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefalutUUID(t *testing.T) {
	gochanUUID = 0
	for i := 0; i < 100; i++ {
		id := defaultUUID()
		assert.Equal(t, i, id)
	}
}
