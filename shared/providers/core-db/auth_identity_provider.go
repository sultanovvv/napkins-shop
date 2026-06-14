package core_db

import (
	"context"
	"database/sql"
	"errors"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

type IAuthIdentityProvider interface {
	Upsert(ctx context.Context, ai entity.AuthIdentity) (*entity.AuthIdentity, error)
	GetByProvider(ctx context.Context, provider, providerUserID string) (*entity.AuthIdentity, error)
}

type authIdentityProvider struct{ db postgres.IBaseProvider }

func NewAuthIdentityProvider(db postgres.IBaseProvider) IAuthIdentityProvider {
	return &authIdentityProvider{db: db}
}

func (p *authIdentityProvider) Upsert(ctx context.Context, ai entity.AuthIdentity) (*entity.AuthIdentity, error) {
	row := models.AuthIdentity{
		UserID:         ai.UserID,
		Provider:       ai.Provider,
		ProviderUserID: ai.ProviderUserID,
	}
	if _, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		On("CONFLICT (provider, provider_user_id) DO UPDATE SET user_id = EXCLUDED.user_id").
		Returning("id, created_at").
		Exec(ctx); err != nil {
		return nil, err
	}
	e := authIdentityToEntity(row)
	return &e, nil
}

func (p *authIdentityProvider) GetByProvider(ctx context.Context, provider, providerUserID string) (*entity.AuthIdentity, error) {
	row := models.AuthIdentity{}
	err := p.db.Conn(ctx).NewSelect().
		Model(&row).
		Where("provider = ?", provider).
		Where("provider_user_id = ?", providerUserID).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := authIdentityToEntity(row)
	return &e, nil
}