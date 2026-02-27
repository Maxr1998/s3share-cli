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

package core

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/maxr1998/s3share-cli/crypto"
	"github.com/maxr1998/s3share-cli/store"
	"github.com/maxr1998/s3share-cli/util"
	"github.com/zeebo/blake3"
)

type UploadInfo struct {
	FileId   string
	FileName string
	Key      string
}

func UploadFile(ctx context.Context, path string) (*UploadInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer util.CloseFileOrExit(file)

	fileId, err := util.GenerateFileId()
	if err != nil {
		return nil, err
	}

	key, err := crypto.GenerateAes256Key()
	if err != nil {
		return nil, err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileName := fileStat.Name()
	fileSize := fileStat.Size()

	// Encrypt filename
	encryptedFileName, err := crypto.EncryptBytesToString([]byte(fileName), key)
	if err != nil {
		return nil, err
	}

	// Encrypt file
	fileCtx, err := crypto.MakeAesCtrContext(key)
	if err != nil {
		return nil, err
	}

	hashers := map[string]hash.Hash{
		"BLAKE3": blake3.New(),
		"MD5":    md5.New(),
	}
	// Go slices aren't covariant…
	hashWriters := make([]io.Writer, 0, len(hashers))
	for _, hasher := range hashers {
		hashWriters = append(hashWriters, hasher)
	}
	description := fmt.Sprintf("Uploading %s", fileName)
	progressReader := util.NewProgressReaderProvider(fileCtx.Encrypt(io.TeeReader(file, io.MultiWriter(hashWriters...))), description, fileSize)
	if err = store.UploadData(ctx, fileId, progressReader, fileSize); err != nil {
		return nil, err
	}

	// Encrypt checksums
	encryptedChecksums := make(map[string]crypto.EncryptedValue, len(hashers))
	for algorithm, hasher := range hashers {
		hasher.Sum(nil)
		encryptedChecksum, err := crypto.EncryptBytesToString(hasher.Sum(nil), key)
		if err != nil {
			return nil, err
		}
		encryptedChecksums[algorithm] = *encryptedChecksum
	}

	// Store metadata
	metadata := store.FileMetadata{
		Name:      *encryptedFileName,
		Checksums: encryptedChecksums,
		Iv:        base64.StdEncoding.EncodeToString(fileCtx.Iv),
		Size:      fileSize,
	}

	if err = store.WriteFileMetadata(ctx, fileId, metadata); err != nil {
		return nil, err
	}

	return &UploadInfo{
		FileId:   fileId,
		FileName: fileName,
		Key:      base64.RawURLEncoding.EncodeToString(key),
	}, nil
}
