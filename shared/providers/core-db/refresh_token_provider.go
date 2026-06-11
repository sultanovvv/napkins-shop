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

type IRefreshTokenProvider interface {
	Create(ctx context.Context, rt entity.RefreshToken) (*entity.RefreshToken, error)
	GetByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id int64, at time.Time, replacedByID *int64) error
	RevokeAllForUser(ctx context.Context, userID int64, at time.Time) error
}

type refreshTokenProvider struct{ db postgres.IBaseProvider }

func NewRefreshTokenProvider(db postgres.IBaseProvider) IRefreshTokenProvider {
	return &refreshTokenProvider{db: db}
}

func (p *refreshTokenProvider) Create(ctx context.Context, rt entity.RefreshToken) (*entity.RefreshToken, error) {
	row := models.RefreshToken{
		UserID:    rt.UserID,
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		UserAgent: rt.UserAgent,
		IP:        rt.IP,
	}
	if _, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		Returning("id, issued_at").
		Exec(ctx); err != nil {
		return nil, err
	}
	e := refreshTokenToEntity(row)
	return &e, nil
}

func (p *refreshTokenProvider) GetByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	row := models.RefreshToken{}
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
	e := refreshTokenToEntity(row)
	return &e, nil
}

func (p *refreshTokenProvider) Revoke(ctx context.Context, id int64, at time.Time, replacedByID *int64) error {
	_, err := p.db.Conn(ctx).NewUpdate().
		Model((*models.RefreshToken)(nil)).
		Set("revoked_at = ?", at).
		Set("replaced_by_id = ?", replacedByID).
		Where("id = ?", id).
		Where("revoked_at IS NULL").
		Exec(ctx)
	return err
}

func (p *refreshTokenProvider) RevokeAllForUser(ctx context.Context, userID int64, at time.Time) error {
	_, err := p.db.Conn(ctx).NewUpdate().
		Model((*models.RefreshToken)(nil)).
		Set("revoked_at = ?", at).
		Where("user_id = ?", userID).
		Where("revoked_at IS NULL").
		Exec(ctx)
	return err
}