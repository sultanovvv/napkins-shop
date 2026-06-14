package entity

import "time"

// User — учётная запись.
//
// PublicID (UUID) уходит в JWT.sub и наружу в API; ID — внутренний BIGSERIAL,
// в кросс-сервисный язык не светится. PasswordHash — nullable: чистый OIDC-юзер
// без локального пароля имеет PasswordHash == nil.
type User struct {
	ID              int64
	PublicID        string
	Email           string
	PasswordHash    *string
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// AuthIdentity — связка "юзер ↔ внешний провайдер". Для локального пароля
// строки в auth_identities нет: признак "имеет локальный пароль" =
// User.PasswordHash != nil.
type AuthIdentity struct {
	ID             int64
	UserID         int64
	Provider       string
	ProviderUserID string
	CreatedAt      time.Time
}