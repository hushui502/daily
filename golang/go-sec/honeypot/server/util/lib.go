package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

func Address2Ip(address string) (ip string) {
	addr := strings.Split(address, ":")
	if len(addr) > 0 {
		ip = addr[0]
	}

	return ip
}

// md5 function
func MD5(s string) (m string) {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// create a sign by key & md5
func MakeSign(t string, key string) (sign string) {
	sign = MD5(fmt.Sprintf("%s%s", t, key))
	return sign
}
