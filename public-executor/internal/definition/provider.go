package definition

import (
	_ "napkins-shop/public-executor/bootstrap/config"
	_ "net/http"
)

//
//func provideS3Provider(client *minio.Client, cfg *config.AppConfig) s3contract.Provider {
//	return s3provider.NewProvider(client, cfg.S3Config.BucketName)
//}
