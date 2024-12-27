package core

import (
	"context"
	"github.com/maxr1998/s3share-cli/store"
	"slices"
	"time"
)

type UploadedFileMetadata struct {
	FileId       string
	Metadata     store.FileMetadata
	ObjectSize   int64
	Exists       bool
	LastModified time.Time
}

func ListUploadedFiles(ctx context.Context) ([]UploadedFileMetadata, error) {
	files, err := store.ListFileMetadata(ctx)
	if err != nil {
		return nil, err
	}

	objects, err := store.ListDataObjects(ctx)
	if err != nil {
		return nil, err
	}

	// Collect S3 metadata
	uploadedFiles := make([]UploadedFileMetadata, 0, len(files))
	for fileId, metadata := range files {
		object, exists := objects[fileId]

		var objectSize int64
		var lastModified time.Time
		if exists {
			objectSize = *object.Size
			lastModified = *object.LastModified
		}

		uploadedFiles = append(uploadedFiles, UploadedFileMetadata{
			FileId:       fileId,
			Metadata:     metadata,
			Exists:       exists,
			ObjectSize:   objectSize,
			LastModified: lastModified,
		})
	}

	// Sort by last modified time
	slices.SortFunc(uploadedFiles, func(a, b UploadedFileMetadata) int {
		return a.LastModified.Compare(b.LastModified)
	})

	return uploadedFiles, nil
}
