package executor

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dibyajyoti-mandal/code-exec-engine/constants"
)

func UploadDoneFile(
	client *s3.Client,
	bucket, key string,
) error {
	body := strings.NewReader("done problem")

	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   body,
	})
	return err
}

func MarkProblemDone() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	return UploadDoneFile(
		client,
		constants.S3Bucket,
		constants.S3DoneKey,
	)
}

func DownloadTestcase(client *s3.Client, bucket, key, localPath string) error {
	out, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}
	defer out.Body.Close()

	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, out.Body)
	return err
}

func FetchSingleTestcase() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	if err := DownloadTestcase(
		client,
		constants.S3Bucket,
		constants.S3Input,
		"/sandbox/input.txt",
	); err != nil {
		return err
	}

	if err := DownloadTestcase(
		client,
		constants.S3Bucket,
		constants.S3Output,
		"/sandbox/expected.txt",
	); err != nil {
		return err
	}

	return nil
}
