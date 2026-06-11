package models

import (
	"time"

	"github.com/uptrace/bun"
)

// EmailVerificationToken — bun-модель таблицы email_verification_tokens.
type EmailVerificationToken struct {
	bun.BaseModel `bun:"table:email_verification_tokens,alias:evt"`

	ID        int64      `bun:"id,pk,autoincrement"`
	UserID    int64      `bun:"user_id,notnull"`
	TokenHash string     `bun:"token_hash,notnull"`
	ExpiresAt time.Time  `bun:"expires_at,notnull"`
	UsedAt    *time.Time `bun:"used_at"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:now()"`
}
