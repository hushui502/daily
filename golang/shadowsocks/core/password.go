package core

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const PasswordLength = 256

type Password [PasswordLength]byte

func init() {
	// 防止生成重复密码
	rand.Seed(time.Now().Unix())
}

// 生成一个256个byte随机组合的密码,用于对内容进行加密解密的密钥
func RandPassword() string {
	// 随机生成一个0-255足成的byte数组
	intArr := rand.Perm(PasswordLength)
	password := &Password{}

	for i, v := range intArr {
		password[i] = byte(v)
		if i == v {
			// 确保不会出现索引和值一致
			return RandPassword()
		}
	}

	return password.String()
}

// 采用base64编码把密码转换为字符串
func (password *Password) String() string {
	return base64.StdEncoding.EncodeToString(password[:])
}

// 解析采用base63编码的字符串获取密码
func ParsePassword(passwordString string) (*Password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(passwordString))
	if err != nil || len(bs) != PasswordLength {
		return nil, errors.New("不合法的密码")
	}

	password := Password{}
	copy(password[:], bs)
	bs = nil

	return &password, nil
}