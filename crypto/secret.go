package crypto

import (
	"crypto/rand"
	"fmt"
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
