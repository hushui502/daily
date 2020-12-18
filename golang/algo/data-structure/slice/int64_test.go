package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt64Slice_Exists(t *testing.T) {
	s := Int64Slice{1, 2, 3}

	assert.True(t, s.Exists(1))
	assert.True(t, s.Exists(2))
	assert.True(t, s.Exists(3))
}

func TestInt64Slice_Sort(t *testing.T) {
	s := Int64Slice{1, 3, 2, 5, 4}
	s.Sort()

	assert.Equal(t, Int64Slice{1, 2, 3, 4, 5}, s)
}

func TestInt64Slice_Search(t *testing.T) {
	s := Int64Slice{1, 2, 3}

	assert.Equal(t, s.Search(1), 0)
	assert.Equal(t, s.Search(2), 1)
	assert.Equal(t, s.Search(3), 2)
}

func TestInt64Slice_Insert(t *testing.T) {
	s := Int64Slice{1, 3, 4}
	s = s.Insert(2)

	assert.Equal(t, Int64Slice{1, 2, 3, 4}, s)

	s = s.Insert(5)
	assert.Equal(t, Int64Slice{1, 2, 3, 4, 5}, s)
}