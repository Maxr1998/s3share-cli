// Copyright © 2024 Maxr1998 <max@maxr1998.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/maxr1998/s3share-cli/conf"
	"io"
)

type EncryptionContext struct {
	Stream cipher.Stream
	Iv     []byte
}

// EncryptedValue represents an encrypted value with an IV, both encoded in base64.
type EncryptedValue struct {
	Value string `json:"value"`
	Iv    string `json:"iv"`
}

// MakeAesCtrContext creates an AES-256 CTR encryption context using the given key.
//
// The returned context contains the stream and IV to be used for encryption.
func MakeAesCtrContext(key []byte) (*EncryptionContext, error) {
	iv, err := generateIv()
	if err != nil {
		return nil, err
	}

	// Create the AES cipher and stream
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)

	return &EncryptionContext{Stream: stream, Iv: iv}, nil
}

// Encrypt encrypts data from the given reader using the encryption context's stream.
func (ctx *EncryptionContext) Encrypt(reader io.Reader) io.Reader {
	return cipher.StreamReader{S: ctx.Stream, R: reader}
}

func (ctx *EncryptionContext) EncryptBytes(data []byte) []byte {
	encrypted := make([]byte, len(data))
	ctx.Stream.XORKeyStream(encrypted, data)
	return encrypted
}

// EncryptBytesToString encrypts the given data with a fresh encryption context and returns the encrypted value.
func EncryptBytesToString(data []byte, key []byte) (*EncryptedValue, error) {
	ctx, err := MakeAesCtrContext(key)
	if err != nil {
		return nil, err
	}
	encrypted := ctx.EncryptBytes(data)
	return &EncryptedValue{
		Value: base64.StdEncoding.EncodeToString(encrypted),
		Iv:    base64.StdEncoding.EncodeToString(ctx.Iv),
	}, nil
}

// generateIv generates an IV to be used for AES-256 encryption.
//
// The part where the counter resides is left empty (zero).
// This IV should never be reused with the same key!
func generateIv() ([]byte, error) {
	iv := make([]byte, conf.IvLength)
	if _, err := rand.Read(iv[:conf.IvSize]); err != nil {
		return nil, err
	}
	return iv, nil
}
