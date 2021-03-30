package api

import (
	"crypto/md5"
	"fmt"
	"honeypot/agent/settings"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	SECRET string
	APIURL string
)

func MD5(s string) string {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func MakeSign(t string, key string) string {
	return MD5(fmt.Sprintf("%s%s", t, key))
}

func init() {
	cfg := settings.Cfg
	sec := cfg.Section("server")
	SECRET = sec.Key("SECRET").MustString("SECRET")
	APIURL = sec.Key("API_URL").MustString("http://127.0.0.1/api/send")
}

func Post(data string) error {
	t := time.Now().Format("2001-01-01 12:32:21")
	hostname, _ := os.Hostname()
	_, err := http.PostForm(APIURL, url.Values{"timestamp": {t}, "secureKey": {MakeSign(t, SECRET)}, "data": {data}, "hostname": {hostname}})
	return err
}
