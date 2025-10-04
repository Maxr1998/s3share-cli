package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/maxr1998/s3share-cli/conf"
)

func AddToHistory(url ShareableUrl) {
	historyFile, err := os.OpenFile(conf.HistoryFileLocation, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return
	}
	defer CloseFileOrExit(historyFile)
	_, _ = fmt.Fprintf(historyFile, "%v\n", url)
}

func ReadHistory() map[string]*ShareableUrl {
	historyFile, err := os.Open(conf.HistoryFileLocation)
	if err != nil {
		return nil
	}
	defer CloseFileOrExit(historyFile)

	urls := make(map[string]*ShareableUrl)
	scanner := bufio.NewScanner(historyFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if url, err := ParseUrl(line); err == nil {
			urls[url.FileId] = url
		} else {
			fmt.Println(line)
		}
	}

	return urls
}

// WriteHistory rewrites the history file with the given urls.
func WriteHistory(urls []*ShareableUrl) {
	historyFile, err := os.OpenFile(conf.HistoryFileLocation, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer CloseFileOrExit(historyFile)

	writer := bufio.NewWriter(historyFile)
	for _, url := range urls {
		_, _ = fmt.Fprintf(writer, "%v\n", url)
	}
}
