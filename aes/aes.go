package aes

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func GenerateKeys() (string, error) {
	// key generation
	key := make([]byte, 16)
	// generating random key
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("Error generating AES key: %v", err)
	}
	// returning the key as a hex string
	return hex.EncodeToString(key), nil
}

func SaveKeyFile(filename string) error {
	// generate keys for AES
	key, err := GenerateKeys()
	// evaluating
	if err != nil {
		return fmt.Errorf("Error to save generate file key: %v", err)
	}
	// Save the key to a file
	err = os.WriteFile(filename, []byte(key), 0600)
	if err != nil {
		return fmt.Errorf("Error saving key file: %v", err)
	}
	// getting the key from the file
	fmt.Println("Key generate successfully")
	return nil
}
