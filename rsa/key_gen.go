package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateRSAPubAndPrivateKey(rsa_priv, rsa_pub string) {
	// generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println("Error to generate the RSA key", err)
	}
	// generate public key
	publicKey := &privateKey.PublicKey
	// getting the bytes
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	// generating the pem file
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	// writing key to file with permission
	err = os.WriteFile(rsa_priv, privateKeyPEM, 0644)
	if err != nil {
		fmt.Println(err)
	}
	// getting the bytes
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	// getting the public bytes
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	// writing key to file with permission
	err = os.WriteFile(rsa_pub, publicKeyPEM, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
