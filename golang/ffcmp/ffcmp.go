package ffcmp

import (
	"runtime"
	"sync"
	"unsafe"
)

// A value pointer is the lite handle to an underlying comparable value.
type Value struct {
	_ [0]func()		// prevent people from accidentally using value type as comparable
	cmpVal interface{}
	gen int64
}

// Get returns the comparable value passed to the Get func
// that had returned v
func (v *Value) Get() interface{} {
	return v.cmpVal
}

var (
	mu sync.Mutex
	valMap = map[interface{}]uintptr{}	// to uintptr(*Value)
)


// Get returns a pointer representing the comparable value cmpVal.
// The returned pointer will be the same for Get(v) and Get(v2)
// if and only if v == v2
func Get(cmpVal interface{}) *Value {
	mu.Lock()
	defer mu.Unlock()

	addr, ok := valMap[cmpVal]
	var v *Value
	if ok {
		v = (*Value)((unsafe.Pointer)(addr))
	} else {
		v = &Value{cmpVal:cmpVal}
		valMap[cmpVal] = uintptr(unsafe.Pointer(v))
	}
	curGen := v.gen + 1
	v.gen = curGen

	if curGen > 1 {
		runtime.SetFinalizer(v, nil)
	}

	runtime.SetFinalizer(v, func(v *Value) {
		mu.Lock()
		defer mu.Unlock()

		if v.gen != curGen {
			// somebody is still using us
			return
		}
		delete(valMap, v.cmpVal)
	})

	return v
}
