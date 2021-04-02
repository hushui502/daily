package misc

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)

	return fmt.Sprintf("%v", h.Sum(nil))
}

func MakeSign(t string, key string) string {
	sign := MD5(fmt.Sprintf("%s%s", t, key))

	return sign
}
