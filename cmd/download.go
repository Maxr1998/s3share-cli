package cmd

import (
	"fmt"
	"github.com/maxr1998/s3share-cli/core"
	"github.com/maxr1998/s3share-cli/util"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [url]…",
	Short: "Download a file",
	Long:  `Downloads and decrypts the file from the given URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("missing URL[s] to download"))
		}

		var urls = make([]*util.ShareableUrl, len(args))

		// Validate the given URls
		for i, arg := range args {
			if url, err := util.ParseUrl(arg); err != nil || url.ServiceHost == "" {
				cobra.CheckErr(fmt.Errorf("invalid URL: %s", arg))
			} else {
				urls[i] = url
			}
		}

		// Download all files sequentially
		ctx := cmd.Context()
		for _, url := range urls {
			downloadFileName, err := core.DownloadFile(ctx, *url)
			cobra.CheckErr(err)

			fmt.Printf("Downloaded %s\n", downloadFileName)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
