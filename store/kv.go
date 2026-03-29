package store

import (
	"context"
	"encoding/json"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/maxr1998/s3share-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var cfAccountId *cloudflare.ResourceContainer
var cfApi *cloudflare.API
var metadataNamespaceId string
var uploadSessionsNamespaceId string

func InitKvClient() {
	var err error
	accountId := viper.GetString("kv.account_id")
	apiToken := viper.GetString("kv.api_token")
	metadataNamespaceId = viper.GetString("kv.namespace_id")

	if accountId == "" || apiToken == "" || metadataNamespaceId == "" {
		println("Missing configuration. Please make sure to set the account ID, API token and namespace ID in your configuration file.")
		os.Exit(1)
	}

	cfAccountId = cloudflare.AccountIdentifier(accountId)
	cfApi, err = cloudflare.NewWithAPIToken(apiToken)
	cobra.CheckErr(err)
}

func InitUploadSessionsNamespaceId() {
	uploadSessionsNamespaceId = viper.GetString("kv.upload_sessions_namespace_id")
	if uploadSessionsNamespaceId == "" {
		println("Missing configuration. Please make sure to set the upload sessions namespace ID in your configuration file.")
		os.Exit(1)
	}
}

// ReadKvData reads the value stored under the given key from the KV store.
func ReadKvData(ctx context.Context, namespaceId string, key string) ([]byte, error) {
	return cfApi.GetWorkersKV(ctx, cfAccountId, cloudflare.GetWorkersKVParams{
		NamespaceID: namespaceId,
		Key:         key,
	})
}

// WriteKvData writes the given value to the KV store under the given key.
func WriteKvData(ctx context.Context, namespaceId string, key string, value []byte) error {
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

// DeleteKvData deletes the KV data stored under the given key.
func DeleteKvData(ctx context.Context, namespaceId string, key string) error {
	response, err := cfApi.DeleteWorkersKVEntry(ctx, cfAccountId, cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: namespaceId,
		Key:         key,
	})
	if err == nil && !response.Success {
		err = util.CollectResponseErrors(response)
	}
	return err
}

// ListKvKeys returns a list of all keys in the KV store.
func ListKvKeys(ctx context.Context, namespaceId string) ([]string, error) {
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

// HasFileMetadata checks if a file with the given file ID exists in the KV store.
func HasFileMetadata(ctx context.Context, fileId string) bool {
	_, err := ReadKvData(ctx, metadataNamespaceId, fileId)
	return err == nil
}

// ReadFileMetadata returns the metadata stored for a file with the given file ID.
func ReadFileMetadata(ctx context.Context, fileId string) (FileMetadata, error) {
	var metadata FileMetadata
	metadataJson, err := ReadKvData(ctx, metadataNamespaceId, fileId)
	if err == nil {
		err = json.Unmarshal(metadataJson, &metadata)
	}
	return metadata, err
}

func WriteFileMetadata(ctx context.Context, fileId string, metadata FileMetadata) error {
	if metadataJson, err := json.Marshal(metadata); err == nil {
		return WriteKvData(ctx, metadataNamespaceId, fileId, metadataJson)
	} else {
		return err
	}
}

func DeleteFileMetadata(ctx context.Context, fileId string) error {
	return DeleteKvData(ctx, metadataNamespaceId, fileId)
}

// ListFileMetadata returns a map with file IDs and their metadata.
func ListFileMetadata(ctx context.Context) (map[string]FileMetadata, error) {
	fileIds, err := ListKvKeys(ctx, metadataNamespaceId)
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

// CreateUploadSession creates a new pending upload session for the given token in the sessions KV namespace.
func CreateUploadSession(ctx context.Context, token string) error {
	return WriteKvData(ctx, uploadSessionsNamespaceId, token, []byte(`{"state":"PENDING"}`))
}

func GetUploadSession(ctx context.Context, token string) (UploadSession, error) {
	var session UploadSession
	sessionJson, err := ReadKvData(ctx, uploadSessionsNamespaceId, token)
	if err == nil {
		err = json.Unmarshal(sessionJson, &session)
	}
	return session, err
}

// ListUploadSessions returns a list of all current upload sessions. TODO: deduplicate with ListFileMetadata
func ListUploadSessions(ctx context.Context) (map[string]UploadSession, error) {
	uploadTokens, err := ListKvKeys(ctx, uploadSessionsNamespaceId)
	if err != nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	// Collect all upload sessions
	results := make([]UploadSession, len(uploadTokens))
	for i, fileId := range uploadTokens {
		g.Go(func() error {
			session, err := GetUploadSession(ctx, fileId)
			if err == nil {
				results[i] = session
			}
			return err
		})
	}

	if err = g.Wait(); err != nil {
		return nil, err
	}

	// Collect results
	sessions := make(map[string]UploadSession, len(uploadTokens))
	for i, fileId := range uploadTokens {
		sessions[fileId] = results[i]
	}

	return sessions, nil
}
