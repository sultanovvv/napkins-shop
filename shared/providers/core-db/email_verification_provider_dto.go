package core_db

import (
	"shared/entity"
	"shared/providers/core-db/models"
)

func emailVerificationToEntity(r models.EmailVerificationToken) entity.EmailVerificationToken {
	return entity.EmailVerificationToken{
		ID:        r.ID,
		UserID:    r.UserID,
		TokenHash: r.TokenHash,
		ExpiresAt: r.ExpiresAt,
		UsedAt:    r.UsedAt,
		CreatedAt: r.CreatedAt,
	}
}
