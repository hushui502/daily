package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// 初始化超时时间为 1 秒
	http.DefaultClient.Timeout = time.Second
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		count := 1
		for {
			select {
			case <-ticker.C:
				// 每隔 5 秒，更新一下超时时间
				http.DefaultClient.Timeout = time.Second * time.Duration(count)
				count++
			}
		}
	}()

	// 不断请求 Google，会触发超时，如果没有超时，说明你已经违法，😄
	for i := 0; i < 100; i++ {
		startTime := time.Now()
		func() {
			resp, err := http.Get("https://www.google.com")
			if err != nil {
				return
			}
			defer resp.Body.Close()
		}()

		// 打印下运行数据，开始时间，超时时间
		fmt.Println(fmt.Sprintf("Run %d:", i+1), "Start:", startTime.Format("15:04:05"),
			"Timeout:", time.Since(startTime))

		// 每隔 1 秒请求一次
		<-time.After(time.Second)
	}
}