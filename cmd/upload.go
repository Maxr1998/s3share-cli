// Copyright © 2024 Maxr1998 <max@maxr1998.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"errors"
	"fmt"
	"github.com/maxr1998/s3share-cli/conf"
	"github.com/maxr1998/s3share-cli/core"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"time"
)

// uploadCmd implements the upload [file]… command
var uploadCmd = &cobra.Command{
	Use:   "upload [file]…",
	Short: "Upload a file",
	Long:  `Encrypts and uploads the specified file to ` + conf.ServiceName + `.`,
	Run: func(cmd *cobra.Command, paths []string) {
		if len(paths) < 1 {
			cobra.CheckErr(fmt.Errorf("missing file[s] to upload"))
		}

		// Verify all given paths exist
		for _, path := range paths {
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				cobra.CheckErr(fmt.Errorf("file %s does not exist", path))
			} else if err != nil {
				cobra.CheckErr(err)
			}
		}

		// Upload all files sequentially
		ctx := cmd.Context()
		for _, path := range paths {
			uploadInfo, err := core.UploadFile(ctx, path, progressReaderProvider)
			cobra.CheckErr(err)

			baseUrl := viper.GetString(conf.ServiceName + ".url")
			fmt.Printf("Uploaded %s: %s/%s#%s\n", uploadInfo.FileName, baseUrl, uploadInfo.FileId, uploadInfo.Key)
		}
	},
}

// progressReaderProvider wraps the provided upstreamReader and returns a new reader
// that tracks the read process through a progress bar that is shown on screen.
func progressReaderProvider(upstreamReader io.Reader, fileName string, fileSize int64) io.Reader {
	desc := fmt.Sprintf("Uploading %s", fileName)
	bar := progressbar.NewOptions64(
		fileSize,
		progressbar.OptionSetDescription(desc),
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

func init() {
	rootCmd.AddCommand(uploadCmd)
}
