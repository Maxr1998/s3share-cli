package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mattn/go-isatty"
	"github.com/maxr1998/s3share-cli/conf"
	"github.com/maxr1998/s3share-cli/core"
	"github.com/maxr1998/s3share-cli/store"
	"github.com/maxr1998/s3share-cli/util"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var isATerminal = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
var warningColor = text.Colors{text.Bold, text.FgRed}

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
		history := util.ReadHistory()

		tableWriter := table.NewWriter()
		tableWriter.SetOutputMirror(cmd.OutOrStdout())
		if isATerminal {
			table.StyleLight.Format.Footer = text.FormatDefault
			tableWriter.SetStyle(table.StyleLight)
		} else {
			table.StyleDefault.Format.Footer = text.FormatDefault
			tableWriter.SetStyle(table.StyleDefault)
		}
		tableWriter.SuppressEmptyColumns()

		tableWriter.AppendHeader(table.Row{"File ID", "Name", "Size", "Checksum", "Last Modified", "Warnings"})
		totalSize := int64(0)
		for _, file := range files {
			// Initialize metadata with defaults
			var decryptedMetadata = &store.DecryptedMetadata{
				Size: file.Metadata.Size,
			}
			if historyUrl := history[file.FileId]; historyUrl != nil {
				if metadata, err := file.Metadata.Decrypt(history[file.FileId].Key); err == nil {
					decryptedMetadata = metadata
				}
			}

			totalSize += decryptedMetadata.Size
			var checksum string
			if decryptedMetadata.Checksum != nil {
				checksum = fmt.Sprintf("%x", decryptedMetadata.Checksum)
			}
			var lastModified string
			if file.Exists {
				lastModified = humanize.Time(file.LastModified)
			} else {
				lastModified = "N/A"
			}
			var warning string
			if !file.Exists {
				warning = warningColor.Sprint("Missing data")
			} else if decryptedMetadata.Size != file.ObjectSize {
				warning = warningColor.Sprint("Size mismatch, data may be corrupt")
			}

			row := table.Row{
				file.FileId,
				decryptedMetadata.Name,
				humanize.IBytes(uint64(decryptedMetadata.Size)),
				checksum,
				lastModified,
				warning,
			}

			tableWriter.AppendRow(row)
		}

		totalNumFiles := fmt.Sprintf("Total %d files", len(files))
		totalSizeStr := humanize.IBytes(uint64(totalSize))
		tableWriter.AppendFooter(
			table.Row{totalNumFiles, "", totalSizeStr, "", ""},
			table.RowConfig{AutoMerge: true},
		)
		tableWriter.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
