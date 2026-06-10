package s3

// Config — параметры S3-совместимого хранилища (MinIO локально, Yandex
// Object Storage в prod). Endpoint — host[:port] без схемы; схема
// выводится из UseSSL.
type Config struct {
	Endpoint   string `mapstructure:"endpoint"    validate:"required"`
	UseSSL     bool   `mapstructure:"use_ssl"`
	Region     string `mapstructure:"region"`
	BucketName string `mapstructure:"bucket_name" validate:"required"`
	AccessKey  string `mapstructure:"access_key"  validate:"required"`
	SecretKey  string `mapstructure:"secret_key"  validate:"required"`
}
