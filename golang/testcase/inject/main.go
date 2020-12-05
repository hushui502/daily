package main

import (
	"fmt"
	"github.com/codegangsta/inject"
	"gopkg.in/macaron.v1"
	"reflect"
)

type A struct {
	Name string
}

type B struct {
	Name string
}

type I interface {
	GetName() string
}

func (b *B) GetName() string {
	return b.Name
}

type C struct {
	Astruct A `inject`
	Bstruct B `inject`
}

type MyFastInvoker func(arg1 A, arg2 I, arg3 string)

// 自己的invoke，这样就可以避免走反射了
func (invoker MyFastInvoker) Invoke(args []interface{}) ([]reflect.Value, error) {
	if a, ok := args[0].(A); ok {
		fmt.Println(a.Name)
	}

	if b, ok := args[1].(I); ok {
		fmt.Println(b.GetName())
	}
	if c, ok := args[2].(string); ok {
		fmt.Println(c)
	}
	return nil, nil
}

func InjectDemo() {
	a := A{Name:"a name"}
	inject1 := inject.New()
	inject1.Map(a)
	inject1.MapTo(&B{Name:"b name"}, (*I)(nil))
	inject1.Set(reflect.TypeOf("string"), reflect.ValueOf("c name"))
	inject1.Invoke(func(arg1 A, arg2 I, arg3 string) {
		fmt.Println(arg1.Name)
		fmt.Println(arg2.GetName())
		fmt.Println(arg3)
	})

	c := C{}
	inject1.Apply(&c)
	fmt.Println(c.Astruct.Name)

	inject2 := inject.New()
	inject2.Map(a)
	inject2.MapTo(&B{Name:"b name"}, (*I)(nil))
	inject2.Set(reflect.TypeOf("string"), reflect.ValueOf("c name"))
	//inject2.Invoke(MyFastInvoker(nil))
}

func main() {
	InjectDemo()

	a := &A{Name:"a name"}
	m := macaron.Classic()
	m.Map(a)
	m.Get("/", func(a *A) string {
		return "hello world " + a.Name
	})
	m.Run()
}
