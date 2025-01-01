package util

import (
	"errors"
	"github.com/spf13/viper"
	"net/url"
	"strings"
)

type ShareableUrl struct {
	ServiceHost string
	FileId      string
	Key         string
}

func (s ShareableUrl) String() string {
	var buffer strings.Builder
	buffer.WriteString("https://")
	buffer.WriteString(s.ServiceHost)
	buffer.WriteString("/")
	buffer.WriteString(s.FileId)
	if s.Key != "" {
		buffer.WriteString("#")
		buffer.WriteString(s.Key)
	}
	return buffer.String()
}

// ParseUrl parses a shareable URL.
// If the host is missing, a default value will be inserted.
func ParseUrl(value string) (*ShareableUrl, error) {
	parsedUrl, err := url.Parse(value)
	if err != nil {
		return nil, err
	}

	host := parsedUrl.Host
	if host == "" {
		host = viper.GetString("service.host")
	}

	fileId := strings.TrimPrefix(parsedUrl.Path, "/")
	if fileId == "" {
		return nil, errors.New("missing file ID")
	} else if strings.Contains(fileId, "/") {
		return nil, errors.New("invalid file ID")
	}

	return &ShareableUrl{
		ServiceHost: host,
		FileId:      fileId,
		Key:         parsedUrl.Fragment,
	}, nil
}
