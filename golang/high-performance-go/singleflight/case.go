package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
	"time"
)

// ==================== DEMO 1 =============================

//type Result string
//
//func find(ctx context.Context, query string) (Result, error) {
//	return Result(fmt.Sprintf("result for %q", query)), nil
//}
//
//func main() {
//	var g singleflight.Group
//	const n = 5
//	wailted := int32(n)
//	done := make(chan struct{})
//	key := "https://weibo.com/1227368500/H3GIgngon"
//
//	// Do
//	for i := 0; i < n; i++ {
//		go func(j int) {
//			v, _, shared := g.Do(key, func() (i interface{}, err error) {
//				ret, err := find(context.Background(), key)
//				return ret, err
//			})
//			if atomic.AddInt32(&wailted, -1) == 0 {
//				close(done)
//			}
//			fmt.Printf("index: %d, val: %v, shared: %v\n", j, v, shared)
//		}(i)
//	}
//
//	select {
//	case <-done:
//	case <-time.After(time.Second):
//		fmt.Println("Do hangs!")
//	}
//
//	// DoChan
//	for i := 0; i < n; i++ {
//		go func(j int) {
//			ch := g.DoChan(key, func() (i interface{}, err error) {
//				ret, err := find(context.Background(), key)
//				return ret, err
//			})
//
//			// create timeout
//			timeout := time.After(time.Second)
//			var ret singleflight.Result
//			select {
//			case <-timeout:
//				fmt.Println("Timeout!")
//				return
//			case ret = <-ch:
//				fmt.Printf("Index: %d, val: %d, sharead: %d\n", j, ret, ret.Val)
//			}
//		}(i)
//	}
//}


// ==================== DEMO 2 =============================
var count int32
func getArticle(id int) (article string, err error) {
	atomic.AddInt32(&count, 1)
	time.Sleep(time.Duration(count) * time.Millisecond)

	return fmt.Sprintf("article: %d", id), nil
}

func singleFlightGetArticle(sg *singleflight.Group, id int) (string, error) {
	v, err, _ := sg.Do(fmt.Sprintf("%d", id), func() (interface{}, error) {
		return getArticle(id)
	})

	return v.(string), err
}

func main() {
	time.AfterFunc(time.Duration(1)*time.Second, func() {
		atomic.AddInt32(&count, -count)
	})

	var (
		wg sync.WaitGroup
		now = time.Now()
		n = 1000
		sg = &singleflight.Group{}
	)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			// res, _ := getArticle(1)			// 同时发起 1000 次请求，耗时: 1.0022831s
			res, _ := singleFlightGetArticle(sg, 1)		// 同时发起 1000 次请求，耗时: 1.5119ms
			if res != "article: 1" {
				panic("err")
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("同时发起 %d 次请求，耗时: %s", n, time.Since(now))
}
