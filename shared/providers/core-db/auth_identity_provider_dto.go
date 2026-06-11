package core_db

import (
	"shared/entity"
	"shared/providers/core-db/models"
)

func authIdentityToEntity(r models.AuthIdentity) entity.AuthIdentity {
	return entity.AuthIdentity{
		ID:             r.ID,
		UserID:         r.UserID,
		Provider:       r.Provider,
		ProviderUserID: r.ProviderUserID,
		CreatedAt:      r.CreatedAt,
	}
}