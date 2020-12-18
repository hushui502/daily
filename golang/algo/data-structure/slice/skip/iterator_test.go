package skip

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterator(t *testing.T) {
	e1 := newMockEntry(1)
	n1 := newNode(e1, 8)

	iter := &iterator{
		n:     n1,
		first: true,
	}

	assert.True(t, iter.Next())
	assert.Equal(t, e1, iter.Value())
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

	e2 := newMockEntry(2)
	n2 := newNode(e2, 8)
	n1.forward[0] = n2

	iter = &iterator{
		n:     n1,
		first: true,
	}

	assert.True(t, iter.Next())
	assert.Equal(t, e1, iter.Value())
	assert.True(t, iter.Next())
	assert.Equal(t, e2, iter.Value())
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

	iter = nilIterator()
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

}
