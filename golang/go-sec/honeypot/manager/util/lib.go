package util

import (
	"crypto/md5"
	"fmt"
)

type Message struct {
	Status  int
	Message string
}

func MakeMd5(src string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(src)))
}

func EncryptPass(passwd string) string {
	return fmt.Sprintf("%s", MakeMd5(passwd)[5:10])
}
