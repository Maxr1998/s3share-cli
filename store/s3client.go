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

package store

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
)

const maxPartSize = int64(128 * 1024 * 1024) // 128 MB

var client *s3.Client
var defaultBucket string

func InitS3Client() {
	uploadUrl := viper.GetString("upload.url")
	bucket := viper.GetString("upload.bucket")
	accessKeyId := viper.GetString("upload.access_key")
	accessKeySecret := viper.GetString("upload.secret_key")

	if uploadUrl == "" || bucket == "" || accessKeyId == "" || accessKeySecret == "" {
		println("Missing configuration. Please make sure to set the endpoint URL, access key and secret key in your configuration file.")
		os.Exit(1)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	cobra.CheckErr(err)

	client = s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(uploadUrl)
	})
	defaultBucket = bucket
}

// UploadData uploads the data from the given reader to the configured S3 bucket.
func UploadData(ctx context.Context, key string, reader io.Reader, fileSize int64) error {
	if fileSize <= maxPartSize {
		return uploadData(ctx, key, reader, fileSize)
	} else {
		return multipartUploadData(ctx, key, reader, fileSize)
	}
}

func uploadData(ctx context.Context, key string, reader io.Reader, fileSize int64) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(defaultBucket),
		Key:           aws.String(key),
		Body:          reader,
		ContentLength: aws.Int64(fileSize),
	})
	return err
}

func multipartUploadData(ctx context.Context, key string, reader io.Reader, fileSize int64) error {
	upload, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(defaultBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	var partLength int64
	var remaining = fileSize
	var completedParts []types.CompletedPart
	var partNumber int32 = 1

	for remaining != 0 {
		if remaining < maxPartSize {
			partLength = remaining
		} else {
			partLength = maxPartSize
		}

		uploadPartResult, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:        upload.Bucket,
			Key:           upload.Key,
			PartNumber:    aws.Int32(partNumber),
			UploadId:      upload.UploadId,
			Body:          io.LimitReader(reader, partLength),
			ContentLength: aws.Int64(partLength),
		})
		if err != nil {
			_, abortErr := client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
				Bucket:   upload.Bucket,
				Key:      upload.Key,
				UploadId: upload.UploadId,
			})
			if abortErr != nil {
				return abortErr
			}
			return err
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       uploadPartResult.ETag,
			PartNumber: aws.Int32(partNumber),
		})
		remaining -= partLength
		partNumber++
	}

	_, err = client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   upload.Bucket,
		Key:      upload.Key,
		UploadId: upload.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func ListDataObjects(ctx context.Context) (map[string]types.Object, error) {
	listObjectResult, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(defaultBucket),
	})
	if err != nil {
		return nil, err
	}

	objects := make(map[string]types.Object)
	for _, object := range listObjectResult.Contents {
		objects[*object.Key] = object
	}

	return objects, nil
}
