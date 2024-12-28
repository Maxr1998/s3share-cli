package util

import (
	"errors"
	"github.com/cloudflare/cloudflare-go"
	"strings"
)

func CollectResponseErrors(response cloudflare.Response) error {
	if len(response.Errors) == 0 {
		return nil
	}
	messages := make([]string, len(response.Errors))
	for i, e := range response.Errors {
		messages[i] = e.Message
	}
	return errors.New(strings.Join(messages, ", "))
}
