package s3gateway

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"

	"napkins-shop/public-executor/configs"
)

// Gateway — обёртка над minio.Client с предзаданным бакетом. Не делает
// public-read sign — бакет настроен на анонимный download через mc.
type Gateway interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Delete(ctx context.Context, key string) error
}

type gateway struct {
	client *minio.Client
	bucket string
}

func NewGateway(cfg *configs.AppConfig, client *minio.Client) Gateway {
	return &gateway{
		client: client,
		bucket: cfg.S3.BucketName,
	}
}

func (g *gateway) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	if g.bucket == "" {
		return fmt.Errorf("s3gateway: bucket not configured")
	}
	_, err := g.client.PutObject(ctx, g.bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (g *gateway) Delete(ctx context.Context, key string) error {
	if g.bucket == "" {
		return fmt.Errorf("s3gateway: bucket not configured")
	}
	return g.client.RemoveObject(ctx, g.bucket, key, minio.RemoveObjectOptions{})
}
