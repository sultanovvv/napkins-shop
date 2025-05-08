package http

import (
	"github.com/minio/minio-go/v7"
	"net/http"

	"github.com/minio/minio-go/v7/pkg/credentials"

	"napkins-shop/public-executor/bootstrap/config"
)

func ProvideDefaultHTTPClient() *http.Client {
	return &http.Client{}
}

func ProvideS3Client(cfg *config.AppConfig) (*minio.Client, error) {
	return minio.New(cfg.S3().Address, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3().AccessKey, cfg.S3().SecretKey, ""),
		Secure: true,
	})
}
