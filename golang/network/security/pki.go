package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)


// an example program to generate a self-signed X.509 certificate for my web site and store it in a.cer file
func main() {

}

func generatePKI() {
	random := rand.Reader

	var privateKey rsa.PrivateKey
	loadKey("private.key", privateKey)

	now := time.Now()
	then := now.Add(60 * 60 * 24 * 365 * 1000 * 1000 * 1000) // one year
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "hu.name",
			Organization: []string{"hu ff"},
		},
		NotBefore:             now,
		NotAfter:              then,
		SubjectKeyId:          []byte{1, 2, 3, 4},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"hu.name", "localhost"},
	}
	derBytes, err := x509.CreateCertificate(random, &template, &template, &privateKey.PublicKey, &privateKey)
	if err != nil {
		return
	}

	certCerFile, _ := os.Create("hu.name.cer")
	pem.Encode(certCerFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certCerFile.Close()

	keyPEMFile, _ := os.Create("private.pem")
	pem.Encode(keyPEMFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(&privateKey)})
	keyPEMFile.Close()
}

func getPKI() {
	certCerFile, _ := os.Open("hu.name.cer")
	derBytes := make([]byte, 1024) // bigger than file
	count, _ := certCerFile.Read(derBytes)
	certCerFile.Close()

	// trim the bytes to actual length in call
	cert, _ := x509.ParseCertificate(derBytes[0:count])

	fmt.Printf("name is %s\n", cert.Subject.CommonName)
	fmt.Printf("not before %s\n", cert.NotBefore.String())
	fmt.Printf("not after %s\n", cert.NotAfter.String())
}

func loadKey(filename string, key interface{}) {
	inFile, _ := os.Open(filename)
	decoder := gob.NewDecoder(inFile)
	_ = decoder.Decode(key)
	inFile.Close()
}


