package iterator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayInt_Iterator(t *testing.T) {
	data := ArrayInt{1, 2, 3}
	iterator := data.Iterator()

	i := 0
	for iterator.HasNext() {
		assert.Equal(t, data[i], iterator.CurrentItem())
		iterator.Next()
		i++
	}
}
