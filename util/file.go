package util

import (
	"crypto/rand"
	"os"

	"github.com/jxskiss/base62"
	"github.com/maxr1998/s3share-cli/conf"
	"github.com/spf13/cobra"
)

func GenerateFileId() (string, error) {
	fileIdBytes := make([]byte, conf.FileIdLength)
	_, err := rand.Read(fileIdBytes)
	if err != nil {
		return "", err
	}
	return base62.EncodeToString(fileIdBytes), nil
}

func GenerateUploadToken() (string, error) {
	tokenBytes := make([]byte, conf.UploadTokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base62.EncodeToString(tokenBytes), nil
}

func CloseFileOrExit(file *os.File) {
	err := file.Close()
	if err != nil {
		cobra.CheckErr(err)
	}
}
