package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		return
	}
	URL, _ := url.Parse(os.Args[1])

	c := &http.Client{
		Timeout:       15 * time.Second,
	}
	fmt.Println(URL.String())
	request, _ := http.NewRequest("GET", URL.String(), nil)
	httpData, _ := c.Do(request)
	fmt.Println("Status code: ", httpData.Status)
	header, _ := httputil.DumpResponse(httpData, false)
	fmt.Println(header)

	contentType := httpData.Header.Get("Content-Type")
	characterSet := strings.SplitAfter(contentType, "charset=")
	if len(characterSet) > 1 {
		fmt.Println("Character Set:", characterSet[1])
	}

	if httpData.ContentLength == -1 {
		fmt.Println("ContentLength is unknown!")
	} else {
		fmt.Println("ContentLength:", httpData.ContentLength)
	}

	length := 0
	var buffer [1024]byte
	r := httpData.Body
	for {
		n, err := r.Read(buffer[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		length = length + n
	}

	fmt.Println("length ", length)
}
