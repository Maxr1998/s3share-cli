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
	"github.com/maxr1998/s3share-cli/conf"
	"github.com/maxr1998/s3share-cli/store"
	"github.com/spf13/cobra"
)

// Command groups
const (
	Management = "management"
)

var rootCmd = &cobra.Command{
	Use:   conf.AppName,
	Short: "Share files via " + conf.ServiceName + ".",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		isManagement := cmd.GroupID == Management

		conf.InitConfig(isManagement)
		if isManagement {
			store.InitS3Client()
			store.InitKvClient()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}

func init() {
	rootCmd.AddGroup(&cobra.Group{
		ID:    Management,
		Title: "Manage files (admin only)",
	})
}
