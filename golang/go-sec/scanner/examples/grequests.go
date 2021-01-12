package main

import (
	"fmt"
	"github.com/levigross/grequests"
	"log"
	"net/url"
)

func main() {
	proxyUrl, err := url.Parse("http://sec.lu:8080")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := grequests.Get("http://mail.163.com/",
		&grequests.RequestOptions{Proxies: map[string]*url.URL{proxyUrl.Scheme: proxyUrl}})
	fmt.Printf("resp: %v, err: %v\n", resp, err)
}
