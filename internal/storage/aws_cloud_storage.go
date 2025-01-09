package storage

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/pkg/utils"
)

type AWSCloudStorage struct {
	Region     string
	BucketName string
}

func NewAWSCloudStorage() *AWSCloudStorage {
	cfg := config.GetConfig()
	return &AWSCloudStorage{
		Region:     cfg.Storage.Cloud.AWS.Region,
		BucketName: cfg.Storage.Cloud.AWS.BucketName,
	}
}

func (a *AWSCloudStorage) Upload(filename string, data []byte) error {
	svc, err := a.getSvc()
	if err != nil {
		return err
	}
	key := a.buildKey(filename)
	if _, err = svc.PutObject(a.createPutObjectInput(key, data)); err != nil {
		return err
	}
	return nil
}

func (a *AWSCloudStorage) Download(filename string) ([]byte, error) {
	return nil, fmt.Errorf("Not Impl")
}

func (a *AWSCloudStorage) Delete(filename string) error {
	return fmt.Errorf("Not Impl")
}

func (a *AWSCloudStorage) getSvc() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(a.Region),
		Credentials: &credentials.Credentials{},
	})
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

func (a *AWSCloudStorage) createPutObjectInput(key string, data []byte) *s3.PutObjectInput {
	return &s3.PutObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(a.BucketName),
		Body:   bytes.NewReader(data),
	}
}

func (a *AWSCloudStorage) buildKey(filename string) string {
	return filename + "_" + utils.GetUnixTimeStamp()
}
