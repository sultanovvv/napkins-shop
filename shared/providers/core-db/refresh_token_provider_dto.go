package core_db

import (
	"shared/entity"
	"shared/providers/core-db/models"
)

func refreshTokenToEntity(r models.RefreshToken) entity.RefreshToken {
	return entity.RefreshToken{
		ID:           r.ID,
		UserID:       r.UserID,
		TokenHash:    r.TokenHash,
		IssuedAt:     r.IssuedAt,
		ExpiresAt:    r.ExpiresAt,
		RevokedAt:    r.RevokedAt,
		ReplacedByID: r.ReplacedByID,
		UserAgent:    r.UserAgent,
		IP:           r.IP,
	}
}