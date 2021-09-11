package generic_demo

import "testing"

func TestListWith3Member(t *testing.T) {
	l := New[int]()
	l.PushFront(10)
	l.PushFront(20)

	if l.Len() != 2 {
		t.Errorf("generic_demo should have 2 items but have %d items in it", l.Len())
	}
}
