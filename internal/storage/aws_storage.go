package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/domain"
)

var _ domain.ArtifactStore = (*AWSCloudStorage)(nil)

type hashingReader struct {
	r      io.Reader
	hasher interface {
		io.Writer
		Sum([]byte) []byte
	}
	n int64
}

func (h *hashingReader) Read(p []byte) (int, error) {
	n, err := h.r.Read(p)
	if n > 0 {
		h.n += int64(n)
		_, _ = h.hasher.Write(p[:n])
	}
	return n, err
}

type AWSCloudStorage struct {
	client     *s3.Client
	bucketName string
}

func NewAWSCloudStorage(cfg *config.Config) (*AWSCloudStorage, error) {
	awsCfg := cfg.Storage.Cloud.AWS
	creds := credentials.NewStaticCredentialsProvider(awsCfg.AccessKeyId, awsCfg.AccessKeySecret, "")

	optFns := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(awsCfg.Region),
		awsconfig.WithCredentialsProvider(creds),
	}

	loaded, err := awsconfig.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	var client *s3.Client
	if awsCfg.Endpoint != "" && awsCfg.Endpoint != "default" {
		client = s3.NewFromConfig(loaded, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(awsCfg.Endpoint)
			o.UsePathStyle = true
		})
	} else {
		client = s3.NewFromConfig(loaded)
	}

	return &AWSCloudStorage{
		client:     client,
		bucketName: awsCfg.BucketName,
	}, nil
}

func (a *AWSCloudStorage) Save(ctx context.Context, key string, r io.Reader) (domain.ObjectInfo, error) {
	hasher := sha256.New()
	hr := &hashingReader{r: readerWithContext(ctx, r), hasher: hasher}

	_, err := a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
		Body:   hr,
	})
	if err != nil {
		return domain.ObjectInfo{}, fmt.Errorf("s3 upload: %w", err)
	}
	return domain.ObjectInfo{
		Key:          key,
		Size:         hr.n,
		SHA256:       hex.EncodeToString(hasher.Sum(nil)),
		LastModified: time.Now().UTC(),
	}, nil
}

func (a *AWSCloudStorage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := a.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 get: %w", err)
	}
	return out.Body, nil
}

func (a *AWSCloudStorage) Delete(ctx context.Context, key string) error {
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (a *AWSCloudStorage) List(ctx context.Context, prefix string) ([]domain.ObjectInfo, error) {
	var out []domain.ObjectInfo
	paginator := s3.NewListObjectsV2Paginator(a.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(a.bucketName),
		Prefix: aws.String(prefix),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			if strings.HasSuffix(key, ".tmp") {
				continue
			}
			info := domain.ObjectInfo{Key: key, Size: aws.ToInt64(obj.Size)}
			if obj.LastModified != nil {
				info.LastModified = obj.LastModified.UTC()
			}
			out = append(out, info)
		}
	}
	return out, nil
}
