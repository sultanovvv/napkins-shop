package models

import (
	"time"

	"github.com/uptrace/bun"
)

// User — bun-модель таблицы users. Внутренний тип persistence; провайдер
// мапит её в entity.User на границе.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID              int64      `bun:"id,pk,autoincrement"`
	PublicID        string     `bun:"public_id,notnull"`
	Email           string     `bun:"email,notnull"`
	PasswordHash    *string    `bun:"password_hash"`
	EmailVerifiedAt *time.Time `bun:"email_verified_at"`
	CreatedAt       time.Time  `bun:"created_at,notnull,default:now()"`
	UpdatedAt       time.Time  `bun:"updated_at,notnull,default:now()"`
}