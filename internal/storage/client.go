package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/CulipBlue/backend_ednic/internal/config"
)

type Client struct {
	minio         *minio.Client
	publicBucket  string
	privateBucket string
}

func NewClient(cfg config.Config) (*Client, error) {
	endpoint := strings.TrimSpace(cfg.ObjectStorageEndpoint)
	accessKey := strings.TrimSpace(cfg.ObjectStorageAccessKey)
	secretKey := strings.TrimSpace(cfg.ObjectStorageSecretKey)
	publicBucket := strings.TrimSpace(cfg.ObjectStoragePublicBucket)
	privateBucket := strings.TrimSpace(cfg.ObjectStoragePrivateBucket)

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return nil, errors.New("object storage endpoint, access key, and secret key are required")
	}

	normalizedEndpoint, useSSL, err := normalizeEndpoint(endpoint, cfg.ObjectStorageUseSSL)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(normalizedEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		minio:         minioClient,
		publicBucket:  publicBucket,
		privateBucket: privateBucket,
	}, nil
}

func (c *Client) Health(ctx context.Context) error {
	if c == nil || c.minio == nil {
		return errors.New("object storage client is not configured")
	}

	if err := c.checkBucket(ctx, c.publicBucket); err != nil {
		return fmt.Errorf("public bucket: %w", err)
	}

	if err := c.checkBucket(ctx, c.privateBucket); err != nil {
		return fmt.Errorf("private bucket: %w", err)
	}

	return nil
}

func (c *Client) Buckets() map[string]string {
	if c == nil {
		return map[string]string{}
	}

	return map[string]string{
		"public":  c.publicBucket,
		"private": c.privateBucket,
	}
}

func (c *Client) checkBucket(ctx context.Context, bucket string) error {
	if strings.TrimSpace(bucket) == "" {
		return errors.New("bucket name is required")
	}

	exists, err := c.minio.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("bucket %q does not exist", bucket)
	}

	return nil
}

func normalizeEndpoint(rawEndpoint string, fallbackUseSSL bool) (string, bool, error) {
	if !strings.Contains(rawEndpoint, "://") {
		return rawEndpoint, fallbackUseSSL, nil
	}

	parsed, err := url.Parse(rawEndpoint)
	if err != nil {
		return "", false, err
	}
	if parsed.Host == "" {
		return "", false, fmt.Errorf("invalid object storage endpoint %q", rawEndpoint)
	}

	return parsed.Host, parsed.Scheme == "https", nil
}
