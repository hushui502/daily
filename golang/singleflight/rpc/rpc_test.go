package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"net/rpc"
	"sync"
	"testing"
)

func TestRpc(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	singleFlight := new(singleflight.Group)
	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		fn := func() (interface{}, error) {
			var reply string
			err = client.Call("Data.GetData", Arg{Caller:i}, &reply)

			return reply, err
		}

		go func(i int) {
			result, _, _ := singleFlight.Do("foo", fn)
			fmt.Printf("caller %d get result %s\n", i, result)
			wg.Done()
		}(i)
	}
}
