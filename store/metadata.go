package store

import (
	"encoding/base64"
	"github.com/maxr1998/s3share-cli/crypto"
)

type FileMetadata struct {
	Name     crypto.EncryptedValue `json:"name"`
	Checksum crypto.EncryptedValue `json:"checksum"`
	Iv       string                `json:"iv"`
	Size     int64                 `json:"size"`
}

type DecryptedMetadata struct {
	Name     string
	Checksum []byte
	Key      []byte
	Iv       []byte
	Size     int64
}

func (metadata *FileMetadata) Decrypt(keyString string) (*DecryptedMetadata, error) {
	key, err := base64.RawURLEncoding.DecodeString(keyString)
	if err != nil {
		return nil, err
	}
	fileName, err := crypto.DecryptValue(metadata.Name, key)
	if err != nil {
		return nil, err
	}
	checksum, err := crypto.DecryptValue(metadata.Checksum, key)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(metadata.Iv)
	if err != nil {
		return nil, err
	}

	return &DecryptedMetadata{
		Name:     string(fileName),
		Checksum: checksum,
		Key:      key,
		Iv:       iv,
		Size:     metadata.Size,
	}, nil
}
