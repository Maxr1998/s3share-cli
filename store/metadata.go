package store

import "github.com/maxr1998/s3share-cli/crypto"

type FileMetadata struct {
	Name crypto.EncryptedValue `json:"name"`
	Iv   string                `json:"iv"`
	Size int64                 `json:"size"`
}
