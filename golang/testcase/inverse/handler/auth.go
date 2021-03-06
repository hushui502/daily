package handler

import (
	"encoding/base64"
	"errors"
	"github.com/panjf2000/goproxy/handler"
	"inverse/config"
	"log"
	"net/http"
	"strings"
)

func (ps *ProxyServer) Auth(rw http.ResponseWriter, req *http.Request) bool {
	var err error
	if config.RuntimeViper.GetBool("server.auth") {
		if ps.Browser, err = ps.auth(rw, req); err != nil {
			return false
		}
		return true
	}
	ps.Browser = "Anonymous"
	return true
}

func (ps *ProxyServer) auth(rw http.ResponseWriter, req *http.Request) (string, error) {
	auth := req.Header.Get("Proxy-Authorization")
	auth = strings.Replace(auth, "Basic", "", 1)

	if auth == "" {
		NeedAuth(rw)
		return "", errors.New("need proxy authorization")
	}

	data, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", errors.New("fail to decoding Proxy-Authorization")
	}

	var user, password string

	userPasswordPair := strings.Split(string(data), ":")
	if len(userPasswordPair) != 2 {
		NeedAuth(rw)
		return "", errors.New("fail to log in")
	}

	user = userPasswordPair[0]
	password = userPasswordPair[1]
	if Verify(user, password) == false {
		NeedAuth(rw)
		return "", errors.New("fail to log in")
	}

	return user, nil
}

func NeedAuth(rw http.ResponseWriter) {
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	defer Client.Close()

	if err != nil {
		log.Printf("fail to get TCP connection of client in auth, %v", err)
	}

	_, _ = Client.Write(HTTP407)

}

func Verify(user, password string) bool {
	if user != "" && password != "" {
		if pass, ok := config.RuntimeViper.GetStringMapString("server.user")[user]; ok && pass == password {
			return true
		}
	}
	return false
}