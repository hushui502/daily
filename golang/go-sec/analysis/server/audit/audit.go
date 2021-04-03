package audit

import (
	"encoding/json"
	"fmt"
	"go-sec/analysis/server/models"
	"go-sec/analysis/server/settings"
	"go-sec/analysis/server/util"
	"io/ioutil"
	"net/http"
)

var (
	EvilIpUrl string
)

func init() {
	sec := settings.Cfg.Section("EVIL_IPS")
	EvilIpUrl = sec.Key("API_URL").MustString("")
}

func PacketAudit(sensorIp string, connInfo models.ConnectionInfo) (err error, result bool, detail models.IplistApi) {
	ips := make([]string, 0)
	ips = append(ips, connInfo.SrcIp, connInfo.DstIp)

	for _, ip := range ips {
		if ip == sensorIp {
			continue
		}
		evilUrl := fmt.Sprintf("%v/api/ip/%v", EvilIpUrl, ip)
		resp, err := http.Get(evilUrl)
		var detail models.IplistApi
		if err == nil {
			ret, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				err = json.Unmarshal(ret, &detail)
				result = detail.Evil
				// util.Log.Debugf("check ip:%v, result: %v, detail: %v", ip, result, detail)
				if result {
					evilConnInfo := models.NewEvilConnectionInfo(sensorIp, connInfo, detail)
					evilConnInfo.Insert()
				}
			}
		}
	}

	return err, result, detail
}

func HttpAudit(sensorIp string, req models.HttpReq) (err error, result bool, evilReq *models.EvilHttpReq) {
	util.Log.Debugf("sensorIp: %v, req: %v", sensorIp, req)
	evilReq = models.NewEvilHttpReq(sensorIp, result, req)
	evilReq.Insert()
	return err, result, evilReq
}

func DnsAudit(sensorIp string, dns models.Dns) (err error, result bool, evilDns *models.EvilDns) {
	util.Log.Debugf("sensorIp: %v, req: %v", sensorIp, dns)
	evilDns = models.NewEvilDns(sensorIp, result, dns)
	err = evilDns.Insert()
	return err, result, evilDns
}
