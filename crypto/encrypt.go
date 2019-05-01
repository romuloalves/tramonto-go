package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// EncryptArtifact encrypts the data with the given secret
func EncryptArtifact(secret string, data []byte) ([]byte, error) {
	// Converts secret to the key
	hash := createHash(secret)

	// Generates the derivated keys
	keys := derivateKeys(hash)

	// Gets the right key to artifacts
	key := keys[1]

	// Creates cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// EncryptConfigFile encrypts the data with the given secret
func EncryptConfigFile(secret string, data []byte) ([]byte, error) {
	// Converts secret to the key
	hash := createHash(secret)

	// Generates the derivated keys
	keys := derivateKeys(hash)

	// Gets the right key to the configuration file
	key := keys[0]

	// Creates cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}
