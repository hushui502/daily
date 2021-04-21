package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func main() {

	/*
		encrypt-decrypt
	*/

	// The GenerateKey method takes in a reader that returns random bits, and
	// the number of bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// The public key is a part of the *crypto.PrivateKey struct
	publicKey := privateKey.PublicKey

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&publicKey,
		[]byte("hello world"),
		nil,
	)
	if err != nil {
		return
	}
	fmt.Println(string(encryptedBytes))

	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OEAPOptions in the end signify that we encrypted the data using OEAP, and that we used
	// SHA256 to hash the input.
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return
	}
	// We get back the original information in the form of bytes, which we
	// the cast to a string and print
	fmt.Println(string(decryptedBytes))


	/*
		signature
	*/

	msg := []byte("verified message")

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator.
	// our private key, the hashing algorithm that we used, and the hash sum of our message
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}

	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(&publicKey, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	}

	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
}
