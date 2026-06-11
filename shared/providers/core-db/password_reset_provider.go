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

type IPasswordResetProvider interface {
	Create(ctx context.Context, prt entity.PasswordResetToken) (*entity.PasswordResetToken, error)
	GetByHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id int64, at time.Time) error
}

type passwordResetProvider struct{ db postgres.IBaseProvider }

func NewPasswordResetProvider(db postgres.IBaseProvider) IPasswordResetProvider {
	return &passwordResetProvider{db: db}
}

func (p *passwordResetProvider) Create(ctx context.Context, prt entity.PasswordResetToken) (*entity.PasswordResetToken, error) {
	row := models.PasswordResetToken{
		UserID:    prt.UserID,
		TokenHash: prt.TokenHash,
		ExpiresAt: prt.ExpiresAt,
	}
	if _, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		Returning("id, created_at").
		Exec(ctx); err != nil {
		return nil, err
	}
	e := passwordResetToEntity(row)
	return &e, nil
}

func (p *passwordResetProvider) GetByHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	row := models.PasswordResetToken{}
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
	e := passwordResetToEntity(row)
	return &e, nil
}

func (p *passwordResetProvider) MarkUsed(ctx context.Context, id int64, at time.Time) error {
	_, err := p.db.Conn(ctx).NewUpdate().
		Model((*models.PasswordResetToken)(nil)).
		Set("used_at = ?", at).
		Where("id = ?", id).
		Where("used_at IS NULL").
		Exec(ctx)
	return err
}