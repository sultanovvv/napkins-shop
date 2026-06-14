package configs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"shared/configs/auth"
	"shared/configs/http_server"
	"shared/configs/postgres"
	"shared/configs/s3"
)

// AppConfig — корневой конфиг сервиса. Секции в shared/, чтобы admin и
// другие сервисы импортировали ровно те же типы.
type AppConfig struct {
	Postgres        *postgres.Config    `mapstructure:"postgres"`
	HTTP            *http_server.Config `mapstructure:"http"`
	S3              *s3.Config          `mapstructure:"s3"`
	Auth            *auth.Config        `mapstructure:"auth"`
	ShutdownTimeout time.Duration       `mapstructure:"shutdown_timeout"`
}

func NewAppConfig() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./public-executor")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
