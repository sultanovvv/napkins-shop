package s3storage

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"shared/configs/s3"
)

// NewClient собирает minio.Client из S3-секции конфига. Подходит для
// локального MinIO и для Yandex Object Storage — оба S3-совместимы.
func NewClient(cfg s3.Config) (*minio.Client, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("s3.endpoint is empty")
	}
	cli, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("minio.New: %w", err)
	}
	return cli, nil
}