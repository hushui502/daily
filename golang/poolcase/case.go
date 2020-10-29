package main

import (
	"fmt"
	"sync"
)

var pool *sync.Pool

type Person struct {
	Name string
}

func initPool() {
	pool = &sync.Pool{New: func() interface{} {
		fmt.Println("create a1 new person")
		return new(Person)
	}}
}

func main() {
	initPool()

	p := pool.Get().(*Person)
	fmt.Println("first get person from pool :", p)

	p.Name = "zhangfei"

	pool.Put(p)
	fmt.Println("second get person from pool :", pool.Get().(*Person))

}
