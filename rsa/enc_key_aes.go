package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func RSAEncryptionWithAES(public_key, key_aes string) {
	// getting the public key PEM
	publicKeyPEM, _ := ioutil.ReadFile(public_key)
	// getting the public key block
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	// parsing the public key
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		fmt.Println(err)
	}
	// then, we read the file
	readFile, _ := ioutil.ReadFile("aes_key.txt")
	plainText := []byte(readFile)
	// get the cipher text
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		fmt.Println(err)
	}
	// now we write the AES to the file
	ioutil.WriteFile("aes_key.txt.enc_rsa", cipherText, 0777)
}
