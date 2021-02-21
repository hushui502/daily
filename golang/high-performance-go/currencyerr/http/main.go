package main

import (
	"fmt"
	"net/http"
)

// 并发中的错误处理
// 可以用一个全局的count来记录
// 这里只适用于对并发的准确度要求不高的场景
// HTTP示例

type Result struct {
	Err error
	Response *http.Response
}

func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					continue
				}
				select {
				case <-done:
					return
				case results <- Result{Err: err, Response: resp}:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)

	errCount := 0
	urls := []string{"www.baidu.com", "www.sougou.com"}
	for result := range checkStatus(done, urls...) {
		if result.Err != nil {
			fmt.Printf("error: %v\n", result.Err)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
