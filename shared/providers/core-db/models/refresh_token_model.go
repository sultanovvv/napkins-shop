package models

import (
	"time"

	"github.com/uptrace/bun"
)

// RefreshToken — bun-модель таблицы refresh_tokens. TokenHash — SHA-256 hex
// от выданного клиенту секрета; сам секрет в БД не хранится.
type RefreshToken struct {
	bun.BaseModel `bun:"table:refresh_tokens,alias:rt"`

	ID           int64      `bun:"id,pk,autoincrement"`
	UserID       int64      `bun:"user_id,notnull"`
	TokenHash    string     `bun:"token_hash,notnull"`
	IssuedAt     time.Time  `bun:"issued_at,notnull,default:now()"`
	ExpiresAt    time.Time  `bun:"expires_at,notnull"`
	RevokedAt    *time.Time `bun:"revoked_at"`
	ReplacedByID *int64     `bun:"replaced_by_id"`
	UserAgent    *string    `bun:"user_agent"`
	IP           *string    `bun:"ip"`
}