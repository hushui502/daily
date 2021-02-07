package fpgo

import "reflect"

func init() {
	MakeNone(&None)
	MakeNone(&NoneInt)
	MakeNone(&NoneString)
}

var None func(fn, slice_or_map interface{}) bool

var NoneInt func(func(value int) bool, []int) bool

var NoneString func(func(value string) bool, []string) bool

func MakeNone(fn interface{}) {
	Maker(fn, none)
}

func none(values []reflect.Value) []reflect.Value {
	fn, col := extractArgs(values)
	var ret bool

	if col.Kind() == reflect.Slice {
		noneSlice(fn, col)
	}
	if col.Kind() == reflect.Map {
		noneMap(fn, col)
	}

	return Valueize(reflect.ValueOf(ret))
}

func noneSlice(fn, s reflect.Value) bool {
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		if ok := callPredicate(fn, v); ok {
			return false
		}
	}
	
	return true
}

func noneMap(fn, m reflect.Value) bool {
	for _, k := range m.MapKeys() {
		v := m.MapIndex(k)
		if ok := callPredicate(fn, v); ok {
			return false
		}
	}
	return true
}