package fpgo

import "reflect"

func init() {
	MakeAny(&Any)
	MakeAny(&AnyInt)
	MakeAny(&AnyString)
}


var Any func(fn, slice_or_map interface{}) bool

var AnyInt func(func(value int) bool, []int) bool

var AnyString func(func(value string) bool, []string) bool

func MakeAny(fn interface{}) {
	Maker(fn, any)
}

func any(values []reflect.Value) []reflect.Value {
	fn, col := extractArgs(values)
	var ret bool
	if col.Kind() == reflect.Map {
		ret = anyMap(fn, col)
	}
	if col.Kind() == reflect.Slice {
		ret = anySlice(fn, col)
	}

	return Valueize(reflect.ValueOf(ret))
}

func anySlice(fn, s reflect.Value) bool {
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		if ok := callPredicate(fn, v); ok {
			return true
		}
	}
	return false
}

func anyMap(fn, m reflect.Value) bool {
	for _, k := range m.MapKeys() {
		v := m.MapIndex(k)
		if ok := callPredicate(fn, v); ok {
			return true
		}
	}

	return false
}

