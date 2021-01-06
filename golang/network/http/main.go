package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	get()
}

// Head
func head() {
	if len(os.Args) != 2 {
		return
	}

	url := os.Args[1]

	response, err := http.Head(url)
	if err != nil {
		return
	}

	fmt.Println(response.Status)
	for k, v := range response.Header {
		fmt.Println(k, ": ",  v)
	}
}

// Get
func get() {

	if len(os.Args) != 2 {
		return
	}

	url := os.Args[1]

	response, err := http.Get(url)
	if err != nil {
		return
	}
	if response.Status != "200 OK" {
		return
	}
	b, _ := httputil.DumpResponse(response, false)
	fmt.Println(string(b))

	fmt.Println("====================")
	contentTypes := response.Header["Content-Type"]

	// usually this is a matter of negotiation between user agent and server
	if !acceptableCharset(contentTypes) {
		return
	}


	var buf [512]byte
	reader := response.Body
	for {
		n, err := reader.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		fmt.Println(string(buf[0:n]))
	}
}

func acceptableCharset(contentTypes []string) bool {
	for _, cType := range contentTypes {
		if strings.Index(cType, "text/html") != -1 {
			return true
		}
	}

	return false
}

// Request
func customReq() {
	url, err := url.Parse(os.Args[1])

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}

	request.Header.Add("Accept-Charset", "UTF-8;q=1,ISO-8859;q=0")
}

// Client
func customClient() {
	url, err := url.Parse(os.Args[1])
	if err != nil {
		return
	}

	client := &http.Client{}

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}

	request.Header.Add("Accept-Charset", "UTF-8;q=1,ISO-8859;q=0")

	response, err := client.Do(request)
	if response.Status != "200 OK" {
		return
	}

	chSet := getCharset(response)
	fmt.Printf("got charset %s\n", chSet)
	if chSet != "UTF-8" {
		return
	}

	var buf [512]byte
	reader := response.Body
	defer reader.Close()
	fmt.Println("got body")
	for {
		n, err := reader.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		fmt.Println(string(buf[0:n]))
	}
}

// return the contentType of response
func getCharset(response *http.Response) string {
	contentType := response.Header.Get("Content-Type")
	if contentType == "" {
		return "UTF-8"
	}

	idx := strings.Index(contentType, "charset:")
	if idx == -1 {
		return "UTF-8"
	}

	return strings.Trim(contentType[idx:], " ")
}


// Proxy
// eg. go run proxyget.go http://xyzproxy.com:8080/ http://www.baidu.com
func proxyGet() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], "http://proxy-host:port ", "http://host:port/page")
		return
	}

	proxyString := os.Args[1]
	proxyUrl, err := url.Parse(proxyString)
	if err != nil {
		return
	}

	rawURL := os.Args[2]
	url, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	transport := &http.Transport{Proxy:http.ProxyURL(proxyUrl)}
	client := &http.Client{Transport:transport}

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}

	dump, _ := httputil.DumpRequest(request, false)
	fmt.Println(string(dump))

	response, err := client.Do(request)
	if err != nil {
		return
	}

	fmt.Println("Read OK")

	if response.Status != "200 OK" {
		return
	}
	fmt.Println("Response OK")

	var buf [512]byte
	reader := response.Body
	defer reader.Close()
	for {
		n, err := reader.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		fmt.Println(string(buf[0:n]))
	}

	os.Exit(1)
}


// Authenticating proxy
const auth = "huhuname:huhupassword"
func authProxy()  {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], "http://proxy-host:port ", "http://host:port/page")
		return
	}

	proxyString := os.Args[1]
	proxyUrl, err := url.Parse(proxyString)
	if err != nil {
		return
	}

	rawURL := os.Args[2]
	url, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	// encode the auth
	basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	transport := &http.Transport{Proxy:http.ProxyURL(proxyUrl)}
	client := &http.Client{Transport:transport}

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}
	request.Header.Add("Proxy-Authorization", basic)

	dump, _ := httputil.DumpRequest(request, false)
	fmt.Println(string(dump))

	response, err := client.Do(request)
	if err != nil {
		return
	}

	fmt.Println("Read OK")

	if response.Status != "200 OK" {
		return
	}
	fmt.Println("Response OK")

	var buf [512]byte
	reader := response.Body
	defer reader.Close()
	for {
		n, err := reader.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		fmt.Println(string(buf[0:n]))
	}

	os.Exit(1)
}


// File server
func fileServer() {
	// deliver files from the directory /var/www
	// fileServer := http.FileServer(http.Dir("var/www"))
	fileServer := http.FileServer(http.Dir("/home/httpd/html"))
	http.Handle("/", fileServer)

	http.HandleFunc("/cgi-bin/printenv", printEnv)

	// register the handler and deliver requests to it
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func printEnv(w http.ResponseWriter, r *http.Request) {
	env := os.Environ()
	w.Write([]byte("<h1>Environment</h1>\n<pre>"))

	for _, v := range env {
		w.Write([]byte(v + "\n"))
	}

	w.Write([]byte("</pre>"))
}

// Server Handler
func serverHandler() {
	myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just return no content-arbitrary headers can be set, arbitrary body
		w.WriteHeader(http.StatusNoContent)
	})

	http.ListenAndServe(":8080", myHandler)
}