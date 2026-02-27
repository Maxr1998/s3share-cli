package core

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"

	"github.com/maxr1998/s3share-cli/api"
	"github.com/maxr1998/s3share-cli/crypto"
	"github.com/maxr1998/s3share-cli/util"
)

func DownloadFile(ctx context.Context, url util.ShareableUrl) (string, error) {
	// Get file metadata
	apiResponse, err := api.GetFileMetadata(url)
	if err != nil {
		return "", err
	}
	metadata, err := apiResponse.Metadata.Decrypt(url.Key)
	if err != nil {
		return "", err
	}
	fileCtx, err := crypto.MakeAesCtrDecryptionContext(metadata.Key, metadata.Iv)
	if err != nil {
		return "", err
	}

	// Download file
	resp, err := httpGetWithContext(ctx, apiResponse.DownloadUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	outputFile, err := CreateOutputFile(metadata.Name)
	if err != nil {
		return "", err
	}
	defer util.CloseFileOrExit(outputFile)

	var hasher hash.Hash
	var expectedChecksum []byte
	for algorithm, checksum := range metadata.Checksums {
		switch algorithm {
		case "MD5":
			hasher = md5.New()
		default:
			continue
		}
		expectedChecksum = checksum
		break
	}
	if hasher == nil {
		return "", fmt.Errorf("no supported checksum found")
	}

	decryptedReader := io.TeeReader(fileCtx.Decrypt(resp.Body), hasher)
	description := fmt.Sprintf("Downloading %s", metadata.Name)
	progressReader := util.NewProgressReaderProvider(decryptedReader, description, metadata.Size)
	if written, err := io.Copy(outputFile, progressReader); err != nil || written != metadata.Size {
		defer func(name string) { _ = os.Rename(name, name+".incomplete") }(outputFile.Name())
		return "", fmt.Errorf("download failed")
	}

	// Verify checksum
	if !bytes.Equal(hasher.Sum(nil), expectedChecksum) {
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
