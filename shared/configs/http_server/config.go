package http_server

// Config — параметры HTTP-сервера. Используется во всех сервисах,
// которые поднимают свой echo.Echo (public-executor, admin, ...).
type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port" validate:"required"`
}
