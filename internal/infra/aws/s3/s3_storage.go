package s3

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	appconfig "github.com/risk-place-angola/backend-risk-place/internal/config"
)

type S3StorageService struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3StorageService(cfg *appconfig.AWSConfig) port.StorageService {
	return &S3StorageService{
		client: s3.NewFromConfig(cfg.AwsConfig),
		bucket: cfg.Bucket,
		region: cfg.Region,
	}
}

func (s *S3StorageService) Upload(ctx context.Context, key string, data io.Reader, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        data,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		slog.Error("failed to upload to S3", "key", key, "error", err)
		return err
	}

	slog.Info("uploaded to S3", "key", key)
	return nil
}

func (s *S3StorageService) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		slog.Error("failed to download from S3", "key", key, "error", err)
		return nil, err
	}

	return result.Body, nil
}

func (s *S3StorageService) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		slog.Error("failed to delete from S3", "key", key, "error", err)
		return err
	}

	slog.Info("deleted from S3", "key", key)
	return nil
}

func (s *S3StorageService) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *S3StorageService) GetURL(key string) string {
	return fmt.Sprintf("/api/v1/storage/%s", key)
}
