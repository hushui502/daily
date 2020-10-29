package main

import (
	"log"
	"sync"
	"time"
)

type resp struct {
	k string
	v string
}

func main() {
	res, _ := fetchData()
	log.Print(res)
}

func rpcwork() resp {
	// do some rpc work
	return resp{}
}

func fetchData() (map[string]string, error) {
	var result = map[string]string{} // result is k -> v
	var keys = []string{"a1", "b", "c"}
	var wg sync.WaitGroup
	var m sync.Mutex
	for i := 0; i < len(keys); i++ {
		wg.Add(1)

		go func() {
			m.Lock()
			defer m.Unlock()
			defer wg.Done()

			// do some rpc
			resp := rpcwork()

			result[resp.k] = resp.v
		}()
	}

	waitTimeout(&wg, time.Millisecond)
	return result, nil
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

//func main() {
//	a1 := new(struct{})
//	b := new(struct{})
//	println(a1, b, a1 == b)
//
//	c := new(struct{})
//	d := new(struct{})
//	fmt.Println(c, d, c == d)
//
//}
