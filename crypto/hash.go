package crypto

import (
	"crypto/md5"
)

// createHash creates a hash from the specified key
func createHash(key string) []byte {
	// Creates the hasher function
	hasher := md5.New()

	// Hashes
	hasher.Write([]byte(key))

	return hasher.Sum(nil)
}
