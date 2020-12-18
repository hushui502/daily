package skip

import "data-structure/common"

// Iterator defines an interface that allows a consumer to iterate
// all result of a query.All values will be visited in-order.
type Iterator interface {
	// Next return a bool indicating if there is future value
	// in the iterator and move the iterator to that value
	Next() bool
	// Value returns a Comparator representing the iterator's current
	// position. If there is no value, this returns nil
	Value() common.Comparator
	// exhaust is a helper method that will iterate this iterator
	// to completion and return a list of resulting Entries in order
	exhaust() common.Comparators
}
