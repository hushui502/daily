package fpgo

import "reflect"

func init() {
	MakeEvery(&Every)
	MakeEvery(&EveryInt)
	MakeEvery(&EveryString)
}

var Every func(fn, slice_or_map interface{}) bool

var EveryInt func(func(value int) bool, []int) bool

var EveryString func(func(value string) bool, []string) bool

func MakeEvery(fn interface{}) {
	Maker(fn, every)
}

func every(values []reflect.Value) []reflect.Value {
	fn, col := extractArgs(values)

	var ret bool
	if col.Kind() == reflect.Map {
		ret = everyMap(fn, col)
	}
	if col.Kind() == reflect.Slice {
		ret = everySlice(fn, col)
	}

	return Valueize(reflect.ValueOf(ret))
}

func everySlice(fn, s reflect.Value) bool {
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		if ok := callPredicate(fn, v, reflect.ValueOf(i)); !ok {
			return false
		}
	}
	return true
}

func everyMap(fn, m reflect.Value) bool {
	for _, k := range m.MapKeys() {
		v := m.MapIndex(k)
		if ok := callPredicate(fn, v); !ok {
			return false
		}
	}
	return true
}