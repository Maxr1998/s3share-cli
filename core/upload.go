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
	"crypto/rand"
	"encoding/base64"
	"github.com/jxskiss/base62"
	"github.com/spf13/cobra"
	"io"
	"os"
	"s3share-cli/conf"
	"s3share-cli/crypto"
	"s3share-cli/store"
)

type ProgressReaderProvider func(upstreamReader io.Reader, fileName string, fileSize int64) io.Reader

type UploadInfo struct {
	FileId   string
	FileName string
	Key      string
	Metadata FileMetadata
}

type FileMetadata struct {
	Name   string `json:"name"`
	NameIv string `json:"name_iv"`
	Iv     string `json:"iv"`
	Size   int64  `json:"size"`
}

func UploadFile(path string, progressReaderProvider ProgressReaderProvider) (*UploadInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer closeFileOrExit(file)

	fileId, err := generateFileId()
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
	nameCtx, err := crypto.MakeAesCtrContext(key)
	if err != nil {
		return nil, err
	}
	encryptedFileName := base64.StdEncoding.EncodeToString(nameCtx.EncryptBytes([]byte(fileName)))

	// Encrypt file
	fileCtx, err := crypto.MakeAesCtrContext(key)
	if err != nil {
		return nil, err
	}

	progressReader := progressReaderProvider(fileCtx.Encrypt(file), fileName, fileSize)
	if err = store.UploadData(fileId, progressReader, fileSize); err != nil {
		return nil, err
	}

	return &UploadInfo{
		FileId:   fileId,
		FileName: fileName,
		Key:      base64.RawURLEncoding.EncodeToString(key),
		Metadata: FileMetadata{
			Name:   encryptedFileName,
			NameIv: base64.StdEncoding.EncodeToString(nameCtx.Iv),
			Iv:     base64.StdEncoding.EncodeToString(fileCtx.Iv),
			Size:   fileSize,
		},
	}, nil
}

func generateFileId() (string, error) {
	fileIdBytes := make([]byte, conf.FileIdLength)
	_, err := rand.Read(fileIdBytes)
	if err != nil {
		return "", err
	}
	return base62.EncodeToString(fileIdBytes), nil
}

func closeFileOrExit(file *os.File) {
	err := file.Close()
	if err != nil {
		cobra.CheckErr(err)
	}
}
