package shadowsocks

import (
	"math/rand"
	"time"
)

const PasswordLength = 256

type Password [PasswordLength]byte

func init() {
	// 防止生成重复密码
	rand.Seed(time.Now().Unix())
}

// 生成一个256个byte随机组合的密码,用于对内容进行加密解密的密钥
func RandPassword() *Password {
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

	return password
}
