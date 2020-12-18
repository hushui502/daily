package skip

import (
	"data-structure/common"
	"github.com/stretchr/testify/mock"
)

type mockEntry uint64

func (me mockEntry) Compare(other common.Comparator) int {
	otherU := other.(mockEntry)
	if me == otherU {
		return 0
	}

	if me > otherU {
		return 1
	}

	return -1
}

func newMockEntry(key uint64) mockEntry {
	return mockEntry(key)
}

type mockIterator struct {
	mock.Mock
}

func (mi *mockIterator) Next() bool {
	args := mi.Called()

	return args.Bool(0)
}

func (mi *mockIterator) Value() common.Comparator {
	args := mi.Called()
	result, ok := args.Get(0).(common.Comparator)
	if !ok {
		return nil
	}

	return result
}

func (mi *mockIterator) exhaust() common.Comparator {
	return nil
}

