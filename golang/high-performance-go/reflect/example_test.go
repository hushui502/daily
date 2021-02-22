package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

func TestKind(t *testing.T) {
	for _, v := range []interface{}{"hello", 11, func() {}} {
		switch v := reflect.ValueOf(v); v.Kind() {
		case reflect.String:
			fmt.Println(v.String())
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Println(v.Int())
		default:
			fmt.Printf("%s", v.Kind())
		}
	}
}

func TestMakeFunc(t *testing.T) {
	swap := func(in []reflect.Value) []reflect.Value {
		return []reflect.Value{in[1], in[0]}
	}

	makeSwap := func(fptr interface{}) {
		fn := reflect.ValueOf(fptr).Elem()

		v := reflect.MakeFunc(fn.Type(), swap)

		fn.Set(v)
	}

	var intSwap func(int, int) (int, int)
	makeSwap(&intSwap)
	fmt.Println(intSwap(1, 2))

	var floatSwap func(float64, float64) (float64, float64)
	makeSwap(&floatSwap)
	fmt.Println(floatSwap(1.0, 2.0))
}

