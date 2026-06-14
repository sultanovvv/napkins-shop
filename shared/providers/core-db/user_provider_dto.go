package core_db

import (
	"shared/entity"
	"shared/providers/core-db/models"
)

func userToEntity(r models.User) entity.User {
	return entity.User{
		ID:              r.ID,
		PublicID:        r.PublicID,
		Email:           r.Email,
		PasswordHash:    r.PasswordHash,
		EmailVerifiedAt: r.EmailVerifiedAt,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}