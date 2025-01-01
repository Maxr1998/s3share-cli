package core

import (
	"context"
	"errors"
	"github.com/maxr1998/s3share-cli/store"
)

func DeleteFile(ctx context.Context, fileId string) error {
	// Check if the file exists before enqueuing it for deletion
	if _, err := store.ReadKvData(ctx, fileId); err != nil {
		return errors.New("file not found")
	}
	if err := store.DeleteDataObject(ctx, fileId); err != nil {
		return err
	}
	if err := store.DeleteKvData(ctx, fileId); err != nil {
		return err
	}
	return nil
}
