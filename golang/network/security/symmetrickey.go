package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/blowfish"
)

func main() {
	hash := md5.New()
	bytes := []byte("hello\n")
	hash.Write(bytes)
	hashValue := hash.Sum(nil)
	hashSize := hash.Size()
	for n := 0; n < hashSize; n += 4 {
		var val uint32
		val = uint32(hashValue[n]) << 24 +
			uint32(hashValue[n+1]) << 16 +
			uint32(hashValue[n+2]) << 8 +
			uint32(hashValue[n+3])

		fmt.Printf("%x ", val)
	}

	blowFish()
}


// ----------------------- Symmetric key encryption------------------------
func blowFish() {
	key := []byte("my key")
	cipher, err := blowfish.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
	}
	src := []byte("hello\n\n\n")
	var enc [512]byte

	cipher.Encrypt(enc[0:], src)

	var dec [8]byte
	cipher.Decrypt(dec[0:], enc[0:])
	result := bytes.NewBuffer(nil)
	result.Write(dec[0:8])
	fmt.Println(string(result.Bytes()))
}