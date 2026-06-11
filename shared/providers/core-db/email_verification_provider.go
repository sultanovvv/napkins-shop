package core_db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

type IEmailVerificationProvider interface {
	Create(ctx context.Context, t entity.EmailVerificationToken) (*entity.EmailVerificationToken, error)
	GetByHash(ctx context.Context, tokenHash string) (*entity.EmailVerificationToken, error)
	MarkUsed(ctx context.Context, id int64, at time.Time) error
}

type emailVerificationProvider struct{ db postgres.IBaseProvider }

func NewEmailVerificationProvider(db postgres.IBaseProvider) IEmailVerificationProvider {
	return &emailVerificationProvider{db: db}
}

func (p *emailVerificationProvider) Create(ctx context.Context, t entity.EmailVerificationToken) (*entity.EmailVerificationToken, error) {
	row := models.EmailVerificationToken{
		UserID:    t.UserID,
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt,
	}
	if _, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		Returning("id, created_at").
		Exec(ctx); err != nil {
		return nil, err
	}
	e := emailVerificationToEntity(row)
	return &e, nil
}

func (p *emailVerificationProvider) GetByHash(ctx context.Context, tokenHash string) (*entity.EmailVerificationToken, error) {
	row := models.EmailVerificationToken{}
	err := p.db.Conn(ctx).NewSelect().
		Model(&row).
		Where("token_hash = ?", tokenHash).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := emailVerificationToEntity(row)
	return &e, nil
}

func (p *emailVerificationProvider) MarkUsed(ctx context.Context, id int64, at time.Time) error {
	_, err := p.db.Conn(ctx).NewUpdate().
		Model((*models.EmailVerificationToken)(nil)).
		Set("used_at = ?", at).
		Where("id = ?", id).
		Where("used_at IS NULL").
		Exec(ctx)
	return err
}
