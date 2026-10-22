package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func RSADecryptionWithAES(private_key string) {
	// getting the private PEM
	privatekeyPEM, _ := ioutil.ReadFile(private_key)
	// getting the block
	privatekeyBlock, _ := pem.Decode(privatekeyPEM)
	// parsing the private key
	privateKey, _ := x509.ParsePKCS1PrivateKey(privatekeyBlock.Bytes)
	// reading the file
	readFile, _ := ioutil.ReadFile("aes_key.txt.enc_rsa")
	cipherText := []byte(readFile)
	// getting the plain text
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	// getting the AES decrypted key
	ioutil.WriteFile("key_aes_dec.txt", plainText, 0777)
}
