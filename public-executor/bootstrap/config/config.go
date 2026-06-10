package config

import (
	"time"

	"github.com/spf13/viper"

	"shared/configs/http_server"
	"shared/configs/postgres"
	"shared/configs/s3"
)

// AppConfig — узкий wrapper над тремя секциями shared-конфигов.
// Содержимое каждой секции живёт в shared/, чтобы admin и другие сервисы
// могли импортировать ровно те же типы.
type AppConfig struct {
	db   postgres.Config
	http http_server.Config
	s3   s3.Config
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		db: postgres.Config{
			Host:            viper.GetString("DATABASE_HOST"),
			Port:            viper.GetInt("DATABASE_PORT"),
			Name:            viper.GetString("DATABASE_NAME"),
			User:            viper.GetString("DATABASE_USER"),
			Password:        viper.GetString("DATABASE_PASSWORD"),
			Debug:           viper.GetBool("DATABASE_DEBUG"),
			MaxOpenConns:    viper.GetInt("DATABASE_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DATABASE_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DATABASE_CONN_MAX_LIFE_TIME") * time.Second,
		},
		http: http_server.Config{
			Host: viper.GetString("HTTP_HOST"),
			Port: viper.GetInt("HTTP_PORT"),
		},
		s3: s3.Config{
			Endpoint:   viper.GetString("S3_ENDPOINT"),
			UseSSL:     viper.GetBool("S3_USE_SSL"),
			Region:     viper.GetString("S3_REGION"),
			BucketName: viper.GetString("S3_BUCKET_NAME"),
			AccessKey:  viper.GetString("S3_ACCESS_KEY"),
			SecretKey:  viper.GetString("S3_SECRET_KEY"),
		},
	}
}

func (a *AppConfig) Database() *postgres.Config { return &a.db }
func (a *AppConfig) HTTP() *http_server.Config  { return &a.http }
func (a *AppConfig) S3() *s3.Config             { return &a.s3 }
