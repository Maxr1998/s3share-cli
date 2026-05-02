package core

import (
	"context"
	"errors"

	"github.com/maxr1998/s3share-cli/store"
)

func DeleteFile(ctx context.Context, fileId string) error {
	// Check if the file exists before enqueuing it for deletion
	if !store.HasFileMetadata(ctx, fileId) {
		return errors.New("file not found")
	}
	if err := store.DeleteDataObject(ctx, fileId); err != nil {
		return err
	}
	if err := store.DeleteFileMetadata(ctx, fileId); err != nil {
		return err
	}
	return nil
}
