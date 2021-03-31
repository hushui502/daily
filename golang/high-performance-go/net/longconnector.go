package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// 长连接
// https://studygolang.com/articles/31094
// https://studygolang.com/articles/12040
// https://golang.org/src/net/http/transport.go#L46
// https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang

var HTTPTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   time.Duration(30) * time.Second, // 连接超时时间
		KeepAlive: time.Duration(60) * time.Second, // 保持长连接的时间
	}).DialContext,
	MaxIdleConns:          500,                             // 最大空闲连接数
	IdleConnTimeout:       time.Duration(60) * time.Second, // 空闲连接的超时时间
	ExpectContinueTimeout: 30 * time.Second,                // 等待服务第一个响应的超时时间
	MaxIdleConnsPerHost:   100,                             // 每个host保持的空闲连接数
}


func main() {
	times := 50
	uri := "http://www.baidu.com"

	start := time.Now()
	client := http.Client{}
	for i := 0; i < times; i++ {
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			panic("Http Req Failed " + err.Error())
		}
		resp, err := client.Do(req)
		if err != nil {
			panic("Http Request Failed " + err.Error())
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)
	}
	fmt.Println("Orig GoNet Short Link", time.Since(start))

	start2 := time.Now()
	client2 := http.Client{Transport: HTTPTransport} // 初始化一个带有transport的http的client
	for i := 0; i < times; i++ {
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			panic("Http Req Failed " + err.Error())
		}
		resp, err := client2.Do(req)
		if err != nil {
			panic("Http Request Failed " + err.Error())
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body) // 如果不及时从请求中获取结果，此连接会占用，其他请求服务复用连接，不管使用数据与否都要读取一下
	}
	fmt.Println("Orig GoNet Long Link", time.Since(start2))
}
