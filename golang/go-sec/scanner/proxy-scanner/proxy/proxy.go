package proxy

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go-sec/scanner/proxy-scanner/log"
	"go-sec/scanner/proxy-scanner/models"
	"go-sec/scanner/proxy-scanner/util"
	"sync"
	"time"
)

var (
	DebugMode = false
	ScanNum = 100
	IpList = "iplist.txt"
	Timeout = 10
)

type CheckProxyFunc func(ip string, port int, protocol string) (isProxy bool, proxyInfo models.ProxyInfo, err error)

var (
	httpProxyFunc CheckProxyFunc = CheckHttpProxy
	sockProxyFunc CheckProxyFunc = CheckSockProxy
)

func Scan(ctx *cli.Context) error {
	if ctx.IsSet("debug") {
		DebugMode = ctx.Bool("debug")
	}

	if DebugMode {
		log.Log.Logger.Level = logrus.DebugLevel
	}

	if ctx.IsSet("timeout") {
		Timeout = ctx.Int("timeout")
	}

	if ctx.IsSet("scan_num") {
		ScanNum = ctx.Int("scan_num")
	}

	if ctx.IsSet("filename") {
		IpList = ctx.String("filename")
	}

	startTime := time.Now()

	proxyAddrList := util.ReadProxyAddr(IpList)
	proxyNum := len(proxyAddrList)
	log.Log.Infof("%v proxies will be check", proxyNum)

	scanBatch := proxyNum / ScanNum
	for i := 0; i < scanBatch; i++ {
		log.Log.Debugf("Scanning %v batches", i+1)
		proxies := proxyAddrList[i*ScanNum:(i+1)*ScanNum]
		CheckProxy(proxies)
	}

	if proxyNum%scanBatch > 0 {
		proxies := proxyAddrList[ScanNum*scanBatch:proxyNum]
		CheckProxy(proxies)
	}

	log.Log.Infof("Scan proxies Done, used time: %v\n", time.Since(startTime))
	models.PrintResult()

	return err
}

func CheckProxy(proxyAddr []util.ProxyAddr) {
	var wg sync.WaitGroup
	wg.Add(len(proxyAddr)*(len(HttpProxyProtocol) + len(SockProxyProtocol)))

	for _, addr := range proxyAddr {
		for _, proto := range HttpProxyProtocol {
			go func(ip string, port int, protocol string) {
				defer wg.Done()
				_ = models.SaveProxies(httpProxyFunc(ip, port, protocol))
			}(addr.Ip, addr.Port, proto)
		}

		for proto := range SockProxyProtocol {
			go func(ip string, port int, protocol string) {
				defer wg.Done()
				_ = models.SaveProxies(sockProxyFunc(ip, port, protocol))
			}(addr.Ip, addr.Port, proto)
		}
	}

	wg.Wait()
}