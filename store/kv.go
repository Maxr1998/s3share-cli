package store

import (
	"context"
	"encoding/json"
	"github.com/cloudflare/cloudflare-go"
	"github.com/maxr1998/s3share-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"os"
)

var cfAccountId *cloudflare.ResourceContainer
var cfApi *cloudflare.API
var namespaceId string

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
}

func ReadKvData(ctx context.Context, key string) ([]byte, error) {
	return cfApi.GetWorkersKV(ctx, cfAccountId, cloudflare.GetWorkersKVParams{
		NamespaceID: namespaceId,
		Key:         key,
	})
}

// ReadFileMetadata returns the metadata stored for a file with the given file ID.
func ReadFileMetadata(ctx context.Context, fileId string) (FileMetadata, error) {
	var metadata FileMetadata
	metadataJson, err := ReadKvData(ctx, fileId)
	if err == nil {
		err = json.Unmarshal(metadataJson, &metadata)
	}
	return metadata, err
}

func WriteKvData(ctx context.Context, key string, value []byte) error {
	response, err := cfApi.WriteWorkersKVEntry(ctx, cfAccountId, cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: namespaceId,
		Key:         key,
		Value:       value,
	})
	if err == nil && !response.Success {
		err = util.CollectResponseErrors(response)
	}
	return err
}

func WriteFileMetadata(ctx context.Context, fileId string, metadata FileMetadata) error {
	if metadataJson, err := json.Marshal(metadata); err == nil {
		return WriteKvData(ctx, fileId, metadataJson)
	} else {
		return err
	}
}

// ListKvKeys returns a list of all keys in the KV store.
func ListKvKeys(ctx context.Context) ([]string, error) {
	entries, err := cfApi.ListWorkersKVKeys(ctx, cfAccountId, cloudflare.ListWorkersKVsParams{
		NamespaceID: namespaceId,
	})
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(entries.Result))
	for _, entry := range entries.Result {
		keys = append(keys, entry.Name)
	}

	return keys, nil
}

// ListFileMetadata returns a map with file IDs and their metadata.
func ListFileMetadata(ctx context.Context) (map[string]FileMetadata, error) {
	fileIds, err := ListKvKeys(ctx)
	if err != nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	// Collect metadata for each file
	results := make([]FileMetadata, len(fileIds))
	for i, fileId := range fileIds {
		g.Go(func() error {
			metadata, err := ReadFileMetadata(ctx, fileId)
			if err == nil {
				results[i] = metadata
			}
			return err
		})
	}
	if err = g.Wait(); err != nil {
		return nil, err
	}

	// Collect results
	files := make(map[string]FileMetadata, len(fileIds))
	for i, fileId := range fileIds {
		files[fileId] = results[i]
	}

	return files, nil
}
