package models

import (
	"time"

	"github.com/uptrace/bun"
)

// PasswordResetToken — bun-модель таблицы password_reset_tokens.
type PasswordResetToken struct {
	bun.BaseModel `bun:"table:password_reset_tokens,alias:prt"`

	ID        int64      `bun:"id,pk,autoincrement"`
	UserID    int64      `bun:"user_id,notnull"`
	TokenHash string     `bun:"token_hash,notnull"`
	ExpiresAt time.Time  `bun:"expires_at,notnull"`
	UsedAt    *time.Time `bun:"used_at"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:now()"`
}