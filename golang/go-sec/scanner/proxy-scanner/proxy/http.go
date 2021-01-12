package proxy

import (
	"fmt"
	"go-sec/scanner/proxy-scanner/log"
	"go-sec/scanner/proxy-scanner/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	HttpProxyProtocol = []string{"http", "https"}
	WebUrl = "http://email.163.com/"
)

func CheckHttpProxy(ip string, port int, protocol string) (isProxy bool, proxyInfo models.ProxyInfo, err error) {
	proxyInfo.Addr = ip
	proxyInfo.Port = port
	proxyInfo.Protocol = protocol

	rawProtocol := fmt.Sprintf("%v://%v:%v", protocol, ip, port)
	proxyUrl, err := url.Parse(rawProtocol)
	if err != nil {
		return false, proxyInfo, err
	}

	transport := &http.Transport{Proxy:http.ProxyURL(proxyUrl)}
	client := &http.Client{Transport:transport, Timeout:time.Duration(Timeout) * time.Second}

	resp, err := client.Get(WebUrl)
	if err != nil {
		return false, proxyInfo, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, proxyInfo, err
		}
		if strings.Contains(string(body), "<title>网易免费邮箱") {
			isProxy = true
		}
	}

	log.Log.Debugf("Checking proxy: %v, isProxy: %v\n", rawProtocol, isProxy)

	return isProxy, proxyInfo, err


}
