package handler

import (
	"io"
	"net"
	"net/http"
)

type ProxyServer struct {
	Travel *http.Transport
	Browser string
}

//func NewProxyServer() *http.Server {
//	if config.RuntimeViper.GetBool("server.cache") {
//		var cachePoolType cache.CachePoolType
//		if config.RuntimeViper.GetString("server.cache_type") == "redis" {
//			cachePoolType = cache.Redis
//		}
//
//	}
//}
//
//func (ps *ProxyServer) ServerHTTP(rw http.ResponseWriter, req *http.Request) {
//	defer func() {
//		if err := recover(); err != nil {
//			rw.WriteHeader(http.StatusInternalServerError)
//		}
//	}()
//
//	if !ps.Auth(rw, req) {
//		return
//	}
//
//	ps.LoadBalancing(req)
//	defer ps.Done(req)
//
//	if req.Method == "CONNECT" {
//		ps.HttpHandler(rw, req)
//	} else if req.Method == "GET" && config.RuntimeViper.GetBool("server.cache") {
//		ps.
//	}
//
//}

//HttpHandler handles http connections.
func (ps *ProxyServer) HttpHandler(rw http.ResponseWriter, req *http.Request) {
	RmProxyHeaders(req)

	resp, err := ps.Travel.RoundTrip(req)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	ClearHeaders(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode) // writes the response status.

	_, err = io.Copy(rw, resp.Body)
	if err != nil && err != io.EOF {
		return
	}
}
func (ps *ProxyServer) HttpsHandler(rw http.ResponseWriter, req *http.Request) {
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	if err != nil {
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}

	Remote, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		http.Error(rw, "Failed", http.StatusBadGateway)
		return
	}

	_, _ = Client.Write(HTTP200)
	go copyRemoteToClient(ps.Browser, Remote, Client)
	go copyRemoteToClient(ps.Browser, Client, Remote)
}

func copyRemoteToClient(User string, Remote, Client net.Conn) {
	defer func() {
		_ = Remote.Close()
		_ = Client.Close()
	}()

	_, err := io.Copy(Remote, Client)
	if err != nil && err != io.EOF {
		return
	}
}
