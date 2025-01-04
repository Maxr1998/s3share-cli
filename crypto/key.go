package crypto

import (
	"crypto/rand"
	"github.com/maxr1998/s3share-cli/conf"
)

// GenerateAes256Key generates a random AES-256 key.
func GenerateAes256Key() ([]byte, error) {
	key := make([]byte, conf.KeyLength)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
