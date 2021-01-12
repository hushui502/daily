package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	proxyUrl, err := url.Parse("http://sec.lu:8080")
	if err != nil {
		log.Fatal(err)
	}
	Transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	httpClient := &http.Client{Transport: Transport, Timeout: time.Second * 3}

	res, err := httpClient.Get("http://email.163.com/")
	if err != nil {
		log.Fatal(err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", resBody)
}
