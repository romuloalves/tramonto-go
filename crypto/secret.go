package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// GenerateSecret generates a new secret to encrypt/decrypt files
func GenerateSecret() (string, error) {
	// Generates random bytes
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	//Formats secret
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// derivateKeys generates three derivated keys from the secret using hkdf
func derivateKeys(secret []byte) [][]byte {
	// Gets the hash function
	hashFn := sha256.New

	// Generates the derivated keys
	hkdf := hkdf.New(hashFn, secret, nil, nil)

	// Reads each one from the reader and returns the response
	var keys [][]byte
	for index := 0; index < 3; index++ {
		key := make([]byte, 16)
		if _, err := io.ReadFull(hkdf, key); err != nil {
			panic(err)
		}

		keys = append(keys, key)
	}

	return keys
}
