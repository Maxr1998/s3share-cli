package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/maxr1998/s3share-cli/conf"
	"github.com/maxr1998/s3share-cli/core"
	"github.com/spf13/cobra"
	"strings"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List uploaded files",
	Long:    `List all files that have been uploaded to ` + conf.ServiceName + `.`,
	GroupID: Management,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cobra.CheckErr(fmt.Errorf("unexpected argument: %s", strings.Join(args, " ")))
		}

		ctx := cmd.Context()
		files, err := core.ListUploadedFiles(ctx)
		cobra.CheckErr(err)

		tableWriter := table.NewWriter()
		tableWriter.SetOutputMirror(cmd.OutOrStdout())
		tableWriter.SuppressEmptyColumns()

		tableWriter.AppendHeader(table.Row{"File ID", "Size", "Last Modified", "Warnings"})
		totalSize := int64(0)
		for _, file := range files {
			totalSize += file.Metadata.Size
			var lastModified string
			if file.Exists {
				lastModified = humanize.Time(file.LastModified)
			} else {
				lastModified = "N/A"
			}
			row := table.Row{
				file.FileId,
				humanize.IBytes(uint64(file.Metadata.Size)),
				lastModified,
			}
			if !file.Exists {
				row = append(row, "Missing data")
			} else if file.Metadata.Size != file.ObjectSize {
				row = append(row, "Size mismatch, data may be corrupt")
			}
			tableWriter.AppendRow(row)
		}
		tableWriter.AppendFooter(table.Row{"TOTAL", humanize.IBytes(uint64(totalSize)), "N/A"})
		tableWriter.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
