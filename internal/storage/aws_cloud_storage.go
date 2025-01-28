package storage

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sboy99/go-vault/config"
)

type AWSCloudStorage struct {
	Region          string
	BucketName      string
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	Token           string
}

func NewAWSCloudStorage() *AWSCloudStorage {
	cfg := config.GetConfig()
	return &AWSCloudStorage{
		Region:          cfg.Storage.Cloud.AWS.Region,
		BucketName:      cfg.Storage.Cloud.AWS.BucketName,
		AccessKeyId:     cfg.Storage.Cloud.AWS.AccessKeyId,
		AccessKeySecret: cfg.Storage.Cloud.AWS.AccessKeySecret,
		Endpoint:        cfg.Storage.Cloud.AWS.Endpoint,
		Token:           "", // TODO: Add token but not required for now
	}
}

func (a *AWSCloudStorage) Upload(filename string, data []byte) error {
	s3, err := a.getS3()
	if err != nil {
		return err
	}
	s3Object := a.createPutObjectInput(filename, data)
	if _, err = s3.PutObject(s3Object); err != nil {
		return err
	}
	return nil
}

func (a *AWSCloudStorage) Download(filename string) ([]byte, error) {
	s3, err := a.getS3()
	if err != nil {
		return nil, err
	}
	s3Object := a.createGetObjectInput(filename)
	result, err := s3.GetObject(s3Object)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)
	return buf.Bytes(), nil
}

func (a *AWSCloudStorage) Delete(filename string) error {
	return fmt.Errorf("Not Impl")
}

func (a *AWSCloudStorage) getS3() (*s3.S3, error) {
	sess, err := a.getSession()
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

func (a *AWSCloudStorage) getSession() (*session.Session, error) {
	creds := credentials.NewStaticCredentials(a.AccessKeyId, a.AccessKeySecret, a.Token)
	if a.Endpoint == "default" {
		return session.NewSession(&aws.Config{
			Region:      aws.String(a.Region),
			Credentials: creds,
		})
	}
	return session.NewSession(&aws.Config{
		Region:           aws.String(a.Region),
		Endpoint:         aws.String(a.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
	})
}

func (a *AWSCloudStorage) createPutObjectInput(key string, data []byte) *s3.PutObjectInput {
	return &s3.PutObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(a.BucketName),
		Body:   bytes.NewReader(data),
	}
}

func (a *AWSCloudStorage) createGetObjectInput(key string) *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(a.BucketName),
	}
}
