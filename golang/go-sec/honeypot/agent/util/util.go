package util

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/elazarl/goproxy"
	"honeypot/agent/vars"
	"io/ioutil"
	"path/filepath"
)

func GetCurDir() string {
	dir, err := filepath.Abs(filepath.Dir("./"))
	if err != nil {
		return ""
	}
	return dir
}

func ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func setCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

func SetCA() error {
	var err error
	caCert, certErr := ReadFile(vars.CaCert)
	caKey, keyErr := ReadFile(vars.CaKey)
	if certErr == nil && keyErr == nil {
		err = setCA(caCert, caKey)
	}
	return err
}
