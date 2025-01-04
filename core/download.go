package core

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/maxr1998/s3share-cli/api"
	"github.com/maxr1998/s3share-cli/crypto"
	"github.com/maxr1998/s3share-cli/util"
	"io"
	"net/http"
	"os"
)

func DownloadFile(ctx context.Context, url util.ShareableUrl) (string, error) {
	// Decode key
	key, err := base64.RawURLEncoding.DecodeString(url.Key)
	if err != nil {
		return "", err
	}

	// Get file metadata
	apiResponse, err := api.GetFileMetadata(url)
	if err != nil {
		return "", err
	}
	fileNameBytes, err := crypto.DecryptValue(apiResponse.Metadata.Name, key)
	if err != nil {
		return "", err
	}
	fileName := string(fileNameBytes)
	checksum, err := crypto.DecryptValue(apiResponse.Metadata.Checksum, key)
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(apiResponse.Metadata.Iv)
	if err != nil {
		return "", err
	}
	fileCtx, err := crypto.MakeAesCtrDecryptionContext(key, iv)
	if err != nil {
		return "", err
	}

	// Download file
	resp, err := httpGetWithContext(ctx, apiResponse.Url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	outputFile, err := CreateOutputFile(fileName)
	if err != nil {
		return "", err
	}
	defer util.CloseFileOrExit(outputFile)

	hash := md5.New()
	decryptedReader := io.TeeReader(fileCtx.Decrypt(resp.Body), hash)
	description := fmt.Sprintf("Downloading %s", fileName)
	progressReader := util.NewProgressReaderProvider(decryptedReader, description, apiResponse.Metadata.Size)
	if _, err = io.Copy(outputFile, progressReader); err != nil {
		return "", err
	}

	// Verify checksum
	if !bytes.Equal(hash.Sum(nil), checksum) {
		return "", fmt.Errorf("checksum mismatch")
	}

	return outputFile.Name(), nil
}

func httpGetWithContext(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func CreateOutputFile(name string) (*os.File, error) {
	var fileName = name
	for i := 0; ; i++ {
		if i > 0 {
			fileName = fmt.Sprintf("%s.%d", name, i)
		}
		if file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644); err == nil {
			return file, nil
		} else if !os.IsExist(err) {
			return nil, err
		}
	}
}
