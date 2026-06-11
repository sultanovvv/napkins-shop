package core_db

import (
	"shared/entity"
	"shared/providers/core-db/models"
)

func passwordResetToEntity(r models.PasswordResetToken) entity.PasswordResetToken {
	return entity.PasswordResetToken{
		ID:        r.ID,
		UserID:    r.UserID,
		TokenHash: r.TokenHash,
		ExpiresAt: r.ExpiresAt,
		UsedAt:    r.UsedAt,
		CreatedAt: r.CreatedAt,
	}
}