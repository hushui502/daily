package main

import (
	"sync"
)

//func main() {
//	//n, err := snowflake.NewNode(1)
//	//if err != nil {
//	//	println(err)
//	//	os.Exit(1)
//	//}
//	//
//	//for i := 0; i < 3; i++ {
//	//	id := n.Generate()
//	//	fmt.Println("id ", id)
//	//	fmt.Println(
//	//		"node", id.Node(),
//	//		"step", id.Step(),
//	//		"time", id.Time(),
//	//		"\n",
//	//		)
//	//}
//
//	//lock()
//
//	//useLock()
//}

var count int

func lock() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	println(count)
}