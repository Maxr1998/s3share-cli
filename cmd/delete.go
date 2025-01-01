package cmd

import (
	"fmt"
	"github.com/maxr1998/s3share-cli/conf"
	"github.com/maxr1998/s3share-cli/core"
	"github.com/maxr1998/s3share-cli/util"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete uploaded files",
	Long:  `Delete a specified file from ` + conf.ServiceName + `.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("missing URL or file ID to delete"))
		}

		url, err := util.ParseUrl(args[0])
		cobra.CheckErr(err)

		ctx := cmd.Context()
		err = core.DeleteFile(ctx, url.FileId)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
