package util

import (
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"time"
)

// NewProgressReaderProvider wraps the provided upstreamReader and returns a new reader
// that tracks the read process through a progress bar that is shown on screen.
func NewProgressReaderProvider(upstreamReader io.Reader, description string, fileSize int64) io.Reader {
	bar := progressbar.NewOptions64(
		fileSize,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowTotalBytes(true),
		progressbar.OptionThrottle(50*time.Millisecond),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetRenderBlankState(true),
	)
	progressReader := progressbar.NewReader(upstreamReader, bar)
	return &progressReader
}
