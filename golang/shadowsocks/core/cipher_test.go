package core

import (
	"math/rand"
	"reflect"
	"testing"
)

const (
	MB = 1024 * 1024
)

func TestCipher(t *testing.T) {
	password := RandPassword()
	t.Log(password)
	p, _ := ParsePassword(password)
	cipher := NewCipher(p)

	org := make([]byte, PasswordLength)
	for i := 0; i < PasswordLength; i++ {
		org[i] = byte(i)
	}

	tmp := make([]byte, PasswordLength)
	copy(tmp, org)
	t.Log(tmp)

	cipher.encode(tmp)
	t.Log(tmp)
	cipher.decode(tmp)
	t.Log(tmp)

	if !reflect.DeepEqual(org, tmp) {
		t.Errorf("编解码存在问题，数据不一致")
	}
}

func BenchmarkEncode(b *testing.B) {
	password := RandPassword()
	p, _ := ParsePassword(password)
	cipher := NewCipher(p)
	bs := make([]byte, MB)
	// 我们主要测试的是编码解码，忽略前面的密钥生成时间
	b.ResetTimer()
	rand.Read(bs)
	cipher.encode(bs)
}

func BenchmarkDecode(b *testing.B) {
	password := RandPassword()
	p, _ := ParsePassword(password)
	cipher := NewCipher(p)
	bs := make([]byte, MB)
	b.ResetTimer()
	rand.Read(bs)
	cipher.decode(bs)
}