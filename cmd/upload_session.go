package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/maxr1998/s3share-cli/crypto"
	"github.com/maxr1998/s3share-cli/store"
	"github.com/maxr1998/s3share-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadSessionCmd = &cobra.Command{
	Use:     "upload-session",
	GroupID: Management,
	Short:   "Manage upload sessions",
}

var uploadSessionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new upload session URL",
	Run: func(cmd *cobra.Command, args []string) {
		store.InitUploadSessionsNamespaceId()

		token, err := util.GenerateUploadToken()
		cobra.CheckErr(err)

		key, err := crypto.GenerateAes256Key()
		cobra.CheckErr(err)

		ctx := cmd.Context()
		err = store.CreateUploadSession(ctx, token)
		cobra.CheckErr(err)

		serviceHost := viper.GetString("service.host")
		keyString := base64.RawURLEncoding.EncodeToString(key)
		fmt.Printf("https://%s/upload/%s#%s\n", serviceHost, token, keyString)
	},
}

func init() {
	uploadSessionCmd.AddCommand(uploadSessionCreateCmd)
	rootCmd.AddCommand(uploadSessionCmd)
}
