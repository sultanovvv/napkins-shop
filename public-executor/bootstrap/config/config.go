package config

import (
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	dbConfig   DBConfig
	httpConfig httpConfig
	S3Config   s3
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		dbConfig: DBConfig{
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
		httpConfig: httpConfig{
			Host: viper.GetString("HTTP_HOST"),
			Port: viper.GetInt("HTTP_PORT"),
		},
		S3Config: s3{
			AccessKey:  viper.GetString("S3_ACCESS_KEY"),
			SecretKey:  viper.GetString("S3_SECRET_KEY"),
			Address:    viper.GetString("S3_ADDRESS"),
			BucketName: viper.GetString("S3_BUCKET_NAME"),
		},
	}
}

type DBConfig struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	Debug           bool
	ConnMaxLifetime time.Duration // В секундах
	MaxOpenConns    int
	MaxIdleConns    int
}

type httpConfig struct {
	Host string
	Port int
}

type s3 struct {
	AccessKey  string
	SecretKey  string
	Address    string
	BucketName string
}

func (a *AppConfig) Database() *DBConfig {
	return &a.dbConfig
}

func (a *AppConfig) HTTP() *httpConfig {
	return &a.httpConfig
}

func (a *AppConfig) S3() *s3 {
	return &a.S3Config
}
