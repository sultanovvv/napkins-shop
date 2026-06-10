package storage

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"napkins-shop/public-executor/bootstrap/config"
)

// NewMinioClient собирает minio.Client из S3-секции конфига. Подходит как
// для локального MinIO, так и для Yandex Object Storage — оба S3-совместимы.
func NewMinioClient(cfg *config.AppConfig) (*minio.Client, error) {
	s3 := cfg.S3()
	if s3.Endpoint == "" {
		return nil, fmt.Errorf("S3_ENDPOINT is empty")
	}

	cli, err := minio.New(s3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3.AccessKey, s3.SecretKey, ""),
		Secure: s3.UseSSL,
		Region: s3.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("minio.New: %w", err)
	}
	return cli, nil
}
