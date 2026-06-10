package s3storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"

	"shared/configs/s3"
)

// IGateway — обёртка над minio.Client с предзаданным бакетом. Не делает
// public-read sign — бакет настроен на анонимный download через mc.
type IGateway interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Delete(ctx context.Context, key string) error
}

type gateway struct {
	client *minio.Client
	bucket string
}

func NewGateway(cfg s3.Config, client *minio.Client) IGateway {
	return &gateway{
		client: client,
		bucket: cfg.BucketName,
	}
}

func (g *gateway) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	if g.bucket == "" {
		return fmt.Errorf("s3storage: bucket not configured")
	}
	_, err := g.client.PutObject(ctx, g.bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (g *gateway) Delete(ctx context.Context, key string) error {
	if g.bucket == "" {
		return fmt.Errorf("s3storage: bucket not configured")
	}
	return g.client.RemoveObject(ctx, g.bucket, key, minio.RemoveObjectOptions{})
}