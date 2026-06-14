package models

import (
	"time"

	"github.com/uptrace/bun"
)

// AuthIdentity — bun-модель таблицы auth_identities. Заполняется только для
// внешних провайдеров (google/yandex/vk/...); для локального пароля
// строки нет.
type AuthIdentity struct {
	bun.BaseModel `bun:"table:auth_identities,alias:ai"`

	ID             int64     `bun:"id,pk,autoincrement"`
	UserID         int64     `bun:"user_id,notnull"`
	Provider       string    `bun:"provider,notnull"`
	ProviderUserID string    `bun:"provider_user_id,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull,default:now()"`
}