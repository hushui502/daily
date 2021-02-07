package fpgo

import "reflect"

func init() {

}

var workers = 6

// Maker takes a function pointer(fn) and implements it with the given reflection-based function implementation
// Internally uses reflect.MakeFunc
// 这里就是结合func type 和 func impl，利用反射来生成一个func
func Maker(fn interface{}, impl func(args []reflect.Value) (results []reflect.Value)) {
	fnV := reflect.ValueOf(fn).Elem()
	fnI := reflect.MakeFunc(fnV.Type(), impl)
	fnV.Set(fnI)
}

// ToI takes a slice and converts it to be []interface
func ToI(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("ToI expects a slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// Valueize takes a number or arguments and returns them as []reflect.Value
func Valueize(values ...interface{}) []reflect.Value {
	ret := make([]reflect.Value, len(values))

	for i := 0; i < len(values); i++ {
		v := values[i]
		if t := reflect.TypeOf(v).String(); t == "reflect.Value" {
			ret[i] = v.(reflect.Value)
		} else {
			ret[i] = reflect.ValueOf(v)
		}
	}

	return ret
}

// SetWorkers sets the number of workers used by the worker pools
func SetWorkers(w int)  {
	workers = w
}

// extractArgs pulls the args from a []reflect.Value and converts as appropriate to underlying types
func extractArgs(values []reflect.Value) (reflect.Value, reflect.Value) {
	fn := interfaceToValue(values[0])
	col := interfaceToValue(values[1])

	return fn, col
}

func interfaceToValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return reflect.ValueOf(v.Interface())
	}
	return v
}

// makeSlice makes a slice of the Output type of the supplied function, and of the specified capacity
func makeSlice(fn reflect.Value, len int) reflect.Value {
	t := reflect.SliceOf(fn.Type().Out(0))
	return reflect.MakeSlice(t, len, len)
}

func makeWorkerChans(t reflect.Type) (chan []reflect.Value, reflect.Value) {
	job := make(chan []reflect.Value)
	res := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, t), 100)

	return job, res
}

func callPredicate(fn reflect.Value, args ...reflect.Value) bool {
	in := fn.Type().NumIn()
	res := fn.Call(args[0:in])

	return res[0].Bool()
}

func callFn(fn reflect.Value, args ...reflect.Value) []reflect.Value {
	in := fn.Type().NumIn()
	res := fn.Call(args[0:in])

	return res
}

