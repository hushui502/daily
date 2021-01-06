package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	genRSAKeys()

	// ---de
	var privateKey rsa.PrivateKey
	loadKey("private.key", &privateKey)

	var publicKey rsa.PublicKey
	loadKey("public.key", &publicKey)
}

// ---------------------- Public key encryption-------------------
func genRSAKeys() {
	reader := rand.Reader
	bitSize := 512
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return
	}

	fmt.Println("Private key primes ", key.Primes[0].String(), key.Primes[1].String())
	fmt.Println("Private key exponent ", key.D.String())

	publicKey := key.PublicKey
	fmt.Println("Public key modulus ", publicKey.N.String())
	fmt.Println("Public key exponent ", publicKey.E)

	saveGobKey("private.key", key)
	saveGobKey("public.key", publicKey)

	savePEMKey("private.pem", key)
}

func saveGobKey(filename string, key interface{}) {
	outFile, _ := os.Create(filename)
	encoder := gob.NewEncoder(outFile)
	_ = encoder.Encode(key)
	outFile.Close()
}

func savePEMKey(filename string, key *rsa.PrivateKey) {
	outFile, _ := os.Create(filename)
	var privateKey = &pem.Block{Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key)}
	pem.Encode(outFile, privateKey)
	outFile.Close()
}

func loadKey(filename string, key interface{}) {
	inFile, _ := os.Open(filename)
	decoder := gob.NewDecoder(inFile)
	_ = decoder.Decode(key)
	inFile.Close()
}