package slice

import "sort"

type Int64Slice []int64

func (s Int64Slice) Len() int {
	return len(s)
}

func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Int64Slice) Search(x int64) int {
	return sort.Search(len(s), func(i int) bool {
		return s[i] >= x
	})
}

func (s Int64Slice) Sort() {
	sort.Sort(s)
}

func (s Int64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Int64Slice) Exists(x int64) bool {
	i := s.Search(x)
	if i == len(s) {
		return false
	}

	// because Search return a value which maybe higher than x
	return s[i] == x
}

// Insert will insert x into the sorted position in this list
// and return a list with the value added. If this slice has
// not been sorted. Insert's behavior is undefined.
func (s Int64Slice) Insert(x int64) Int64Slice {
	i := s.Search(x)
	if i == len(s) {
		return append(s, x)
	}

	if s[i] == x {
		return s
	}

	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = x

	return s
}