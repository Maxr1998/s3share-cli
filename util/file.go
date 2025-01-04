package util

import (
	"crypto/rand"
	"github.com/jxskiss/base62"
	"github.com/maxr1998/s3share-cli/conf"
	"github.com/spf13/cobra"
	"os"
)

func GenerateFileId() (string, error) {
	fileIdBytes := make([]byte, conf.FileIdLength)
	_, err := rand.Read(fileIdBytes)
	if err != nil {
		return "", err
	}
	return base62.EncodeToString(fileIdBytes), nil
}

func CloseFileOrExit(file *os.File) {
	err := file.Close()
	if err != nil {
		cobra.CheckErr(err)
	}
}
