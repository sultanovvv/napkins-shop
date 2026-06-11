package auth

import "time"

// Config — параметры подсистемы авторизации. Используется в любом сервисе,
// который выпускает или валидирует токены (public-executor, future admin).
//
// JWTSecret — симметричный ключ для HS256. Когда придёт смена на RS256, тут
// появится KeyID + PEM-приватник; интерфейс TokenSigner спрячет различие.
//
// Argon2 — параметры хэширования пароля. Значения по умолчанию заданы в
// usecases/auth/password.go (NewArgon2Hasher), сюда виперу попадают только
// если переопределены в config.yaml.
type Config struct {
	JWTSecret        string        `mapstructure:"jwt_secret"         validate:"required"`
	AccessTokenTTL   time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL  time.Duration `mapstructure:"refresh_token_ttl"`
	PasswordResetTTL time.Duration `mapstructure:"password_reset_ttl"`
	EmailVerifyTTL   time.Duration `mapstructure:"email_verify_ttl"`

	// GuestCartTTL — срок жизни HttpOnly-cookie gct. Сам guest_token_hash
	// в БД не имеет TTL: чистим строки carts через будущий cleanup-job.
	GuestCartTTL time.Duration `mapstructure:"guest_cart_ttl"`

	Argon2 Argon2Config `mapstructure:"argon2"`
}

type Argon2Config struct {
	Memory      uint32 `mapstructure:"memory"`      // KiB
	Iterations  uint32 `mapstructure:"iterations"`
	Parallelism uint8  `mapstructure:"parallelism"`
	SaltLen     uint32 `mapstructure:"salt_len"`
	KeyLen      uint32 `mapstructure:"key_len"`
}