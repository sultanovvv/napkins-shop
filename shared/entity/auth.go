package entity

import "time"

// RefreshToken — серверная запись о выданном refresh-секрете.
//
// В БД хранится TokenHash (SHA-256 от opaque random secret); сам секрет
// возвращается клиенту один раз при выпуске и больше нигде не сохраняется.
// ReplacedByID собирает аудит-цепочку rotation: при /refresh старый помечается
// revoked_at + replaced_by_id указывает на свежевыпущенный.
type RefreshToken struct {
	ID           int64
	UserID       int64
	TokenHash    string
	IssuedAt     time.Time
	ExpiresAt    time.Time
	RevokedAt    *time.Time
	ReplacedByID *int64
	UserAgent    *string
	IP           *string
}

// TokenPair — то, что выдаём клиенту после успешного login/refresh.
// AccessToken — JWT; RefreshTokenSecret — сырой 32-байтный секрет,
// уходит в HttpOnly-cookie на стороне handler'а.
type TokenPair struct {
	AccessToken        string
	AccessExpiresAt    time.Time
	RefreshTokenSecret string
	RefreshExpiresAt   time.Time
}

// PasswordResetToken — одноразовый токен восстановления пароля. В БД хранится
// только TokenHash; сырой Secret уходит в email (когда появится транспорт).
type PasswordResetToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}