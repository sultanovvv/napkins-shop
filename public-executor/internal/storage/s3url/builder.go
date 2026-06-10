package s3url

import (
	"fmt"
	"strings"

	"napkins-shop/public-executor/bootstrap/config"
)

// Builder превращает S3-ключ объекта в публичный URL вида
// {http|https}://{endpoint}/{bucket}/{key}. Подходит для path-style доступа,
// который поддерживают и MinIO, и Yandex Object Storage.
type Builder interface {
	URL(key string) string
}

type builder struct {
	scheme   string
	endpoint string
	bucket   string
}

func NewBuilder(cfg *config.AppConfig) Builder {
	s3 := cfg.S3()
	scheme := "http"
	if s3.UseSSL {
		scheme = "https"
	}
	return &builder{
		scheme:   scheme,
		endpoint: strings.TrimSuffix(s3.Endpoint, "/"),
		bucket:   s3.BucketName,
	}
}

func (b *builder) URL(key string) string {
	if key == "" || b.bucket == "" || b.endpoint == "" {
		return ""
	}
	return fmt.Sprintf("%s://%s/%s/%s", b.scheme, b.endpoint, b.bucket, strings.TrimPrefix(key, "/"))
}
