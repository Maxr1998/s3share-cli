package util

import (
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
