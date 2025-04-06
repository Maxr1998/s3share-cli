package util

import (
	"fmt"
	"github.com/maxr1998/s3share-cli/conf"
	"os"
)

func AddToHistory(url ShareableUrl) {
	historyFile, err := os.OpenFile(conf.HistoryFileLocation, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer CloseFileOrExit(historyFile)
	_, _ = fmt.Fprintf(historyFile, "%v\n", url)
}
