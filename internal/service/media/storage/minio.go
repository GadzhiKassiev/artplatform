package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type MinIOClient struct {
	client *minio.Client
	bucket string
}

func NewMinIO(ctx context.Context, cfg MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("bucket check: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("make bucket: %w", err)
		}
	}

	return &MinIOClient{client: client, bucket: cfg.Bucket}, nil
}

func (m *MinIOClient) Upload(ctx context.Context, key string, content []byte, contentType string) error {
	_, err := m.client.PutObject(ctx, m.bucket, key,
		bytes.NewReader(content), int64(len(content)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func (m *MinIOClient) Download(ctx context.Context, key string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	return io.ReadAll(obj)
}

func (m *MinIOClient) PresignedGetURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucket, key, ttl, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (m *MinIOClient) Delete(ctx context.Context, key string) error {
	return m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
}
