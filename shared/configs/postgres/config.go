package postgres

import "time"

// Config — параметры подключения к Postgres. Используется в любом
// сервисе (public-executor, future admin), который ходит в core-db.
// Теги mapstructure позволяют viper-у развернуть YAML/dotenv в эту структуру.
type Config struct {
	Host            string        `mapstructure:"host"            validate:"required"`
	Port            int           `mapstructure:"port"            validate:"required"`
	Name            string        `mapstructure:"name"            validate:"required"`
	User            string        `mapstructure:"user"            validate:"required"`
	Password        string        `mapstructure:"password"        validate:"required"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	Debug           bool          `mapstructure:"debug"`
}
