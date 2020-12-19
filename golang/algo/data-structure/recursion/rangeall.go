package recursion

import "fmt"

type RangeType struct {
	value []interface{}
}

func NewRangeArray(n int) *RangeType {
	return &RangeType{
		make([]interface{}, n),
	}
}

func (slice *RangeType) RangeAll(start int) {
	len := len(slice.value)
	if start == len-1 {
		fmt.Println(slice.value)
	}

	for i := start; i < len; i++ {
		if i == start || slice.value[i] != slice.value[start] {
			slice.value[i], slice.value[start] = slice.value[start], slice.value[i]
			slice.RangeAll(i+1)
			slice.value[i], slice.value[start] = slice.value[start], slice.value[i]
		}
	}
}