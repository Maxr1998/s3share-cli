package api

import (
	"encoding/json"
	"fmt"
	"github.com/maxr1998/s3share-cli/store"
	"github.com/maxr1998/s3share-cli/util"
	"net/http"
)

type Response struct {
	FileId      string             `json:"file_id"`
	Metadata    store.FileMetadata `json:"metadata"`
	DownloadUrl string             `json:"url"`
}

func GetFileMetadata(url util.ShareableUrl) (*Response, error) {
	response, err := http.Get(fmt.Sprintf("https://%s/download?file=%s", url.ServiceHost, url.FileId))
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse Response
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
