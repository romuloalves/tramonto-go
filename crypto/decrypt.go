package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

// DecryptArtifact decrypts the data with the given secret
func DecryptArtifact(secret string, data []byte) ([]byte, error) {
	// Creates the hash
	hash := createHash(secret)

	// Generates the derivated keys
	keys := derivateKeys(hash)

	// Gets the right key to artifacts
	key := keys[1]

	// Creates the cipher
	block, err := aes.NewCipher(key)
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

// DecryptConfigFile decrypts the data with the given secret
func DecryptConfigFile(secret string, data []byte) ([]byte, error) {
	// Creates the hash
	hash := createHash(secret)

	// Generates the derivated keys
	keys := derivateKeys(hash)

	// Gets the right key to the configuration file
	key := keys[0]

	// Creates the cipher
	block, err := aes.NewCipher(key)
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
