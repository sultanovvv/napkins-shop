package s3storage

import (
	"fmt"
	"strings"

	"shared/configs/s3"
)

// IURLBuilder превращает S3-ключ объекта в публичный URL вида
// {http|https}://{endpoint}/{bucket}/{key}. Подходит для path-style доступа,
// который поддерживают и MinIO, и Yandex Object Storage.
type IURLBuilder interface {
	URL(key string) string
}

type urlBuilder struct {
	scheme   string
	endpoint string
	bucket   string
}

func NewURLBuilder(cfg s3.Config) IURLBuilder {
	scheme := "http"
	if cfg.UseSSL {
		scheme = "https"
	}
	return &urlBuilder{
		scheme:   scheme,
		endpoint: strings.TrimSuffix(cfg.Endpoint, "/"),
		bucket:   cfg.BucketName,
	}
}

func (b *urlBuilder) URL(key string) string {
	if key == "" || b.bucket == "" || b.endpoint == "" {
		return ""
	}
	return fmt.Sprintf("%s://%s/%s/%s", b.scheme, b.endpoint, b.bucket, strings.TrimPrefix(key, "/"))
}