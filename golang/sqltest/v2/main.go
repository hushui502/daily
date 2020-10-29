package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var usage = `Usage: %s [options]
Options are:
    -n number     Number of requests to perform
    -c concurrency  Number of multiple requests to make at a time
    -t timeout      Seconds to max. wait for each response
    -m method       Method name
`

var (
	number int
	concurrency int
	timeout int
	method string
	url string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}

	flag.IntVar(&number, "n", 1000, "")
	flag.IntVar(&concurrency, "c", 10, "")
	flag.IntVar(&timeout, "t", 1, "")
	flag.StringVar(&method, "m", "GET", "")

	if flag.NArg() != 1 {
		exit("Invalid url")
	}

	method = strings.ToUpper(method)
	url = flag.Args()[0]

	if method != "GET" {
		exit("Invalid method")
	}

	if number < 1 || concurrency < 1 {
		exit("-n and -c cannot be smaller than 1.")
	}

	if number < concurrency {
		exit("-n cannot be less than -c.")
	}
}

type benchmark struct {
	number int
	concurrency int
	timeout int
	method string
	url string
	duration chan time.Duration
	start time.Time
	end time.Time
}

func (b *benchmark) run() {
	b.duration = make(chan time.Duration, b.number)
	b.start = time.Now()
	b.runWorks()
	b.end = time.Now()

	b.report()
}

func (b *benchmark) runWorks() {
	var wg sync.WaitGroup

	wg.Add(b.concurrency)

	for i := 0; i < b.concurrency; i++ {
		go func() {
			defer wg.Done()
			b.runWork(b.number / b.concurrency)
		}()
	}

	wg.Wait()
	close(b.duration)
}

func (b *benchmark) runWork(num int) {
	client := &http.Client{
		Timeout:time.Duration(b.timeout) * time.Second,
	}

	for i := 0; i < num; i++ {
		b.request(client)
	}
}

func (b *benchmark) request(client *http.Client) {
	req, err := http.NewRequest(b.method, b.url, nil)

	if err != nil {
		log.Fatal(req)
	}

	start := time.Now()
	client.Do(req)
	end := time.Now()

	b.duration <- end.Sub(start)
}

func (b *benchmark) report() {
	sum := 0.0
	num := float64(len(b.duration))

	for duration := range b.duration {
		sum += duration.Seconds()
	}

	rps := int(num / b.end.Sub(b.start).Seconds())
	tpr := sum / num * 1000

	fmt.Printf("rps: %d [#/sec]\n", rps)
	fmt.Printf("tpr: %.3f [ms]\n", tpr)
}

func exit(msg string) {
	flag.Usage()
	fmt.Fprintln(os.Stderr, "\n[ERROR] " + msg)
	os.Exit(1)
}

func main() {
	b := benchmark{
		number:      number,
		concurrency: concurrency,
		timeout:     timeout,
		method:      method,
		url:         url,
	}

	b.run()
}