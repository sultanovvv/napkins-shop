package configs

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"shared/configs/http_server"
	"shared/configs/postgres"
)

type AppConfig struct {
	Postgres        *postgres.Config   `mapstructure:"postgres"`
	HTTP            http_server.Config `mapstructure:"http"`
	ShutdownTimeout time.Duration      `mapstructure:"shutdown_timeout"`
}

func NewAppConfig() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./backoffice-executor")
	v.AddConfigPath(".")
	v.SetDefault("http.host", "127.0.0.1")
	v.SetDefault("http.port", 9889)
	v.SetDefault("shutdown_timeout", "10s")

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if cfg.HTTP.Port < 1 || cfg.HTTP.Port > 65535 {
		return nil, fmt.Errorf("http.port must be between 1 and 65535")
	}
	if cfg.ShutdownTimeout <= 0 {
		return nil, fmt.Errorf("shutdown_timeout must be positive")
	}
	if cfg.Postgres == nil {
		return nil, fmt.Errorf("postgres config is required: copy backoffice-executor/config.dist.yaml to backoffice-executor/config.yaml")
	}
	if cfg.Postgres.Host == "" || cfg.Postgres.Name == "" || cfg.Postgres.User == "" {
		return nil, fmt.Errorf("postgres.host, postgres.name and postgres.user are required")
	}
	if cfg.Postgres.Port < 1 || cfg.Postgres.Port > 65535 {
		return nil, fmt.Errorf("postgres.port must be between 1 and 65535")
	}

	return &cfg, nil
}
