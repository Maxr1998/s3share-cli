package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"
)

type DecryptionContext struct {
	Stream cipher.Stream
	Iv     []byte
}

// MakeAesCtrDecryptionContext creates an AES-256 CTR decryption context using the given key and IV.
func MakeAesCtrDecryptionContext(key []byte, iv []byte) (*DecryptionContext, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)

	return &DecryptionContext{Stream: stream, Iv: iv}, nil
}

// Decrypt decrypts data from the given reader using the decryption context's stream.
func (ctx *DecryptionContext) Decrypt(reader io.Reader) io.Reader {
	return cipher.StreamReader{S: ctx.Stream, R: reader}
}

func (ctx *DecryptionContext) DecryptBytes(data []byte) []byte {
	decrypted := make([]byte, len(data))
	ctx.Stream.XORKeyStream(decrypted, data)
	return decrypted
}

// DecryptValue decrypts the given data with a fresh decryption context and returns the decrypted value.
func DecryptValue(value EncryptedValue, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(value.Value)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(value.Iv)
	if err != nil {
		return nil, err
	}
	ctx, err := MakeAesCtrDecryptionContext(key, iv)
	if err != nil {
		return nil, err
	}
	return ctx.DecryptBytes(data), nil
}
