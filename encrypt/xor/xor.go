package xor 

func XorBytes(data, key []byte) []byte {
	// masked
	masked := make([]byte, len(data))
	// key length
	keyLen := len(key)
	// looping
	for i, b := range data {
		masked[i] = b ^ key[i % keyLen]
	}
	// returned the converted code
	return masked
}