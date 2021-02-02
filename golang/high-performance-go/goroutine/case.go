//package main
//
//import (
//	"context"
//	"fmt"
//	"net/http"
//	"time"
//)
//
//func main() {
//	req, err := http.NewRequest("GET", "http://www.baidu.com", nil)
//	if err != nil {
//		return
//	}
//
//	ctx, cancel := context.WithTimeout(req.Context(), 2 * time.Second)
//	defer cancel()
//
//	req = req.WithContext(ctx)
//	client := http.DefaultClient
//	resp, err := client.Do(req)
//	if err != nil {
//		return
//	}
//
//	fmt.Printf("%v\n", resp.StatusCode)
//}

package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan bool, 2)
	ch <- true
	ch <- true
	close(ch)

	var wg *sync.WaitGroup
	for v := range ch {
		fmt.Println(v) // called twice
	}

	wg.Wait()
	close(ch)
	return

}