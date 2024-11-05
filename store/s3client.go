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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
)

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

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	cobra.CheckErr(err)

	client = s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(uploadUrl)
	})
	defaultBucket = bucket
}

func UploadData(key string, reader io.Reader, fileSize int64) error {
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        &defaultBucket,
		Key:           &key,
		Body:          reader,
		ContentLength: &fileSize,
	})
	return err
}
