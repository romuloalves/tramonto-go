package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

// Decrypt decrypts the data with the given secret
func Decrypt(secret string, data []byte) ([]byte, error) {
	// Creates the hash
	hash := createHash(secret)

	// Creates the cipher
	block, err := aes.NewCipher(hash)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Reads the nonce size
	nonceSize := gcm.NonceSize()

	// Splits nonce and real data
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	// Reads data
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
