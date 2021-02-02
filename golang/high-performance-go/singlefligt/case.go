package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync/atomic"
	"time"
)

type Result string

func find(ctx context.Context, query string) (Result, error) {
	return Result(fmt.Sprintf("result for %q", query)), nil
}

func main() {
	var g singleflight.Group
	const n = 5
	wailted := int32(n)
	done := make(chan struct{})
	key := "https://weibo.com/1227368500/H3GIgngon"


	// Do
	//for i := 0; i < n; i++ {
	//	go func(j int) {
	//		v, _, shared := g.Do(key, func() (i interface{}, err error) {
	//			ret, err := find(context.Background(), key)
	//			return ret, err
	//		})
	//		if atomic.AddInt32(&wailted, -1) == 0 {
	//			close(done)
	//		}
	//		fmt.Printf("index: %d, val: %v, shared: %v\n", j, v, shared)
	//	}(i)
	//}
	//
	//select {
	//case <-done:
	//case <-time.After(time.Second):
	//	fmt.Println("Do hangs!")
	//}


	// DoChan
	for i := 0; i < n; i++ {
		go func(j int) {
			ch := g.DoChan(key, func() (i interface{}, err error) {
				ret, err := find(context.Background(), key)
				return ret, err
			})

			// create timeout
			timeout := time.After(time.Second)
			var ret singleflight.Result
			select {
			case <-timeout:
				fmt.Println("Timeout!")
				return
			case ret = <-ch:
				fmt.Printf("Index: %d, val: %d, sharead: %d\n", j, ret, ret.Val)
			}
		}(i)
	}
}
