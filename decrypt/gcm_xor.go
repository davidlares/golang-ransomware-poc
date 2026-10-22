package decrypt

import (
	"crypto/cipher"
	"fmt"
	"io"

	"crypto/aes"
	"io/ioutil"
	"os"
	"strings"
)

const (
	CHUNK_SIZE = 0x100000
)

func DecryptWithXorAndAES(filename, keyFile string) error {
	// reading the keyFile
	key, _ := ioutil.ReadFile(keyFile)
	encData, _ := os.Open(filename)
	// closing
	defer encData.Close()
	// getting the raw file
	outFilename := strings.TrimSuffix(filename, ".enc")
	outFile, _ := os.Create(outFilename)
	// closing
	defer outFile.Close()
	// getting the block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating AES")
	}
	// getting the GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error to create AES-GCM")
	}
	// nonce size
	nonceSize := gcm.NonceSize()
	for {
		nonce := make([]byte, nonceSize)
		_, err := io.ReadFull(encData, nonce)
		if err == io.EOF {
			break
		}
		// then
		if err != nil {
			fmt.Println("Error reading nonce")
		}
		// getting the cipher
		cipherText := make([]byte, CHUNK_SIZE+gcm.Overhead())
		n, err := encData.Read(cipherText)
		if err != nil && err != io.EOF {
			fmt.Println("Error to read the cipherText", err)
		}
		// chunking
		cipherText = cipherText[:n]
		// getting the plain text by chunks
		plainText, err := gcm.Open(nil, nonce, cipherText, nil)
		if err != nil {
			fmt.Println("Error in decrypting", err)
		}
		// then
		_, err = outFile.Write(plainText)
		if err != nil {
			fmt.Println("Error creating file", err)
		}
		// and
		if err == io.EOF {
			break
		}
		// deleting enc file
		err = os.Remove(filename)
		if err != nil {
			fmt.Println("Error removing file")
		}
	}
	return nil
}
