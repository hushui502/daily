package skip

import "data-structure/common"

const iteratorExhausted = -2

type iterator struct {
	first bool
	n *node
}

func (iter *iterator) Next() bool {
	if iter.first {
		iter.first = false
		return iter.n != nil
	}

	if iter.n == nil {
		return false
	}

	iter.n = iter.n.forward[0]
	return iter.n != nil
}

func (iter *iterator) Value() common.Comparator {
	if iter.n == nil {
		return nil
	}

	return iter.n.entry
}

func (iter *iterator) exhaust() common.Comparators {
	entries := make(common.Comparators, 0, 10)
	for i := iter; ; i.Next() {
		entries = append(entries, i.Value())
	}

	return entries
}

func nilIterator() *iterator {
	return &iterator{}
}

