package store

import (
	"encoding/base64"
	"errors"

	"github.com/maxr1998/s3share-cli/crypto"
)

type FileMetadata struct {
	Name      crypto.EncryptedValue            `json:"name"`
	Checksum  *crypto.EncryptedValue           `json:"checksum,omitempty"`
	Checksums map[string]crypto.EncryptedValue `json:"checksums,omitempty"`
	Iv        string                           `json:"iv"`
	Size      int64                            `json:"size"`
}

type DecryptedMetadata struct {
	Name      string
	Checksums map[string][]byte
	Key       []byte
	Iv        []byte
	Size      int64
}

var ErrNoChecksums = errors.New("no checksums found")

func (metadata *FileMetadata) Decrypt(keyString string) (*DecryptedMetadata, error) {
	key, err := base64.RawURLEncoding.DecodeString(keyString)
	if err != nil {
		return nil, err
	}
	fileName, err := crypto.DecryptValue(metadata.Name, key)
	if err != nil {
		return nil, err
	}
	checksums := make(map[string][]byte)
	if metadata.Checksum != nil {
		checksum, err := crypto.DecryptValue(*metadata.Checksum, key)
		if err != nil {
			return nil, err
		}
		checksums["MD5"] = checksum
	}
	for algorithm, checksum := range metadata.Checksums {
		checksums[algorithm], err = crypto.DecryptValue(checksum, key)
		if err != nil {
			return nil, err
		}
	}
	if len(checksums) == 0 {
		return nil, ErrNoChecksums
	}

	iv, err := base64.StdEncoding.DecodeString(metadata.Iv)
	if err != nil {
		return nil, err
	}

	return &DecryptedMetadata{
		Name:      string(fileName),
		Checksums: checksums,
		Key:       key,
		Iv:        iv,
		Size:      metadata.Size,
	}, nil
}
