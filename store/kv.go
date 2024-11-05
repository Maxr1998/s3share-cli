package store

import (
	"context"
	"encoding/json"
	"github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfAccountId *cloudflare.ResourceContainer
var cfApi *cloudflare.API
var namespaceId string

var cfCtx context.Context

func InitKvClient() {
	var err error
	accountId := viper.GetString("kv.account_id")
	apiToken := viper.GetString("kv.api_token")
	namespaceId = viper.GetString("kv.namespace_id")

	if accountId == "" || apiToken == "" || namespaceId == "" {
		println("Missing configuration. Please make sure to set the account ID, API token and namespace ID in your configuration file.")
		os.Exit(1)
	}

	cfAccountId = cloudflare.AccountIdentifier(accountId)
	cfApi, err = cloudflare.NewWithAPIToken(apiToken)
	cobra.CheckErr(err)

	cfCtx = context.Background()
}

func WriteKvData(key string, value []byte) error {
	_, err := cfApi.WriteWorkersKVEntry(cfCtx, cfAccountId, cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: namespaceId,
		Key:         key,
		Value:       value,
	})
	return err
}

func WriteFileMetadata(fileId string, metadata FileMetadata) error {
	if metadataJson, err := json.Marshal(metadata); err == nil {
		return WriteKvData(fileId, metadataJson)
	} else {
		return err
	}
}
