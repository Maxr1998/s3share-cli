package cmd

import (
	"fmt"
	"github.com/maxr1998/s3share-cli/conf"
	"github.com/maxr1998/s3share-cli/core"
	"github.com/maxr1998/s3share-cli/util"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete uploaded files",
	Long:    `Delete a specified file from ` + conf.ServiceName + `.`,
	GroupID: Management,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("missing URL or file ID to delete"))
		}

		var urls = make([]*util.ShareableUrl, 0, len(args))
		for _, arg := range args {
			url, err := util.ParseUrl(arg)
			cobra.CheckErr(err)
			urls = append(urls, url)
		}

		ctx := cmd.Context()
		for _, url := range urls {
			err := core.DeleteFile(ctx, url.FileId)
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
