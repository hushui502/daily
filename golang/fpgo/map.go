package fpgo

import "reflect"

func init() {
	MakeMap(&Map)
	MakeMap(&MapString)
	MakeMap(&MapInt)
	MakePMap(&MapP)
	MakePMap(&MapPString)
}

// Map func(func(A) C, []A) []C
var Map func(interface{}, interface{}) []interface{}

var MapString func(func(string) string, []string) []string

var MapInt func(func(int) int, []int) []int

var MapP func(interface{}, interface{}, ...int) []interface{}

var MapPString func(func(string) string, []string, ...int) []string

func MakeMap(fn interface{}) {
	Maker(fn , mapImpl)
}

func MakePMap(fn interface{}) {
	Maker(fn, mapPImpl)
}

func mapImpl(values []reflect.Value) []reflect.Value {
	fn, col := extractArgs(values)

	ret := makeSlice(fn, col.Len())

	if col.Kind() == reflect.Slice {
		ret = mapSlice(fn, col)
	}

	return []reflect.Value{ret}
}

func mapSlice(fn, col reflect.Value) reflect.Value {
	ret := makeSlice(fn, col.Len())
	for i := 0; i < col.Len(); i++ {
		e := col.Index(i)
		r := fn.Call([]reflect.Value{e})
		ret.Index(i).Set(r[0])
	}

	return ret
}

func mapWorker(fn reflect.Value, job chan []reflect.Value, res reflect.Value) {
	for {
		v, ok := <-job
		if !ok {
			break
		}
		if len(v) > 0 {
			r := fn.Call(v)
			res.Send(r[0])
		}
	}
}

func mapPImpl(values []reflect.Value) []reflect.Value {
	fn, col := extractArgs(values)

	workers := 1
	if len(values) == 3 {
		if l := values[2].Len(); l == 1 {
			workers = int(values[2].Index(0).Int())
		}
	}

	t := fn.Type().Out(0)
	job, res := makeWorkerChans(t)

	ret := makeSlice(fn, col.Len())

	for i := 1; i <= workers; i++ {
		go mapWorker(fn, job, res)
	}

	if col.Kind() == reflect.Slice {
		mapPMap(job, col)
	}

	close(job)

	for i := 0; i < col.Len(); i++ {
		v, ok := res.Recv()
		if !ok {
			break
		}
		ret.Index(i).Set(v)
	}

	return []reflect.Value{ret}
}

func mapPSlice(job chan []reflect.Value, col reflect.Value) {
	for i := 0; i < col.Len(); i++ {
		e := col.Index(i)
		job <- []reflect.Value{e}
	}
}

func mapPMap(job chan []reflect.Value, col reflect.Value) {
	for _, k := range col.MapKeys() {
		v := col.MapIndex(k)
		job <- []reflect.Value{v, k}
	}
}

func refMapMap(m map[string]int, fn func(string, int) string) []string {
	ret := make([]string, 0, len(m))
	for k, v := range m {
		ret = append(ret, fn(k, v))
	}

	return ret
}