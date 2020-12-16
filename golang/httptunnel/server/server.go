package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

const (
	readTimeout = 100
	keyLen = 64
)

type proxy struct {
	C chan proxyPacket
	key string
	conn net.Conn
}

type proxyPacket struct {
	c http.ResponseWriter
	r *http.Request
	done chan bool
}

func NewProxy(key, destAddr string) (p *proxy, err error) {
	p = &proxy{C: make(chan proxyPacket), key: key}
	log.Println("Attempting connect ", destAddr)
	p.conn, err = net.Dial("tcp", destAddr)
	if err != nil {
		return
	}

	p.conn.SetReadDeadline(time.Now().Add(time.Microsecond * readTimeout))
	log.Println("ResponseWriter ", destAddr)

	return
}

func (p *proxy) handle(pp proxyPacket) {
	// read from the request and write to the ResponseWriter
	_, err := io.Copy(p.conn, pp.r.Body)
	pp.r.Body.Close()
	if err == io.EOF {
		p.conn = nil
		log.Println("eof ", p.key)
		return
	}

	// read out of the buffer and write it to conn
	pp.c.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(pp.c, p.conn)

	pp.done <- true
}

var queue = make(chan proxyPacket)
var createQueue = make(chan *proxy)

func handler(w http.ResponseWriter, r *http.Request) {
	pp := proxyPacket{w, r, make(chan bool)}
	queue <- pp

	// wait until done before returning
	<-pp.done
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	// read destAddr
	destAddr, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, "Could not found read destAddr",
			http.StatusInternalServerError)

		return
	}

	key := genKey()

	p, err := NewProxy(key, string(destAddr))
	if err != nil {
		http.Error(w, "Could not connect",
			http.StatusInternalServerError)

		return
	}

	createQueue <- p

	w.Write([]byte(key))
}

func proxyMuxer() {
	proxyMap := make(map[string]*proxy)
	for {
		select {
		case pp := <-queue:
			key := make([]byte, keyLen)
			// read key
			n, err := pp.r.Body.Read(key)
			if n != keyLen || err != nil {
				log.Println("Couldn't read key", key)
				continue
			}
			// find proxy
			p, ok := proxyMap[string(key)]
			if !ok {
				log.Println("Couldn't find proxy", key)
				continue
			}
			// handle
			p.handle(pp)
		case p := <-createQueue:
			proxyMap[p.key] = p
		}
	}
}

var httpAddr = flag.String("http", ":8888", "http listen address")

func main() {
	flag.Parse()

	go proxyMuxer()

	http.HandleFunc("/", handler)
	http.HandleFunc("/create", createHandler)
	http.ListenAndServe(*httpAddr, nil)
}

func genKey() string {
	key := make([]byte, keyLen)
	for i := 0; i < keyLen; i++ {
		key[i] = byte(rand.Int())
	}

	return string(key)
}


























