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

type IUserProvider interface {
	Create(ctx context.Context, u entity.User) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByPublicID(ctx context.Context, publicID string) (*entity.User, error)
	UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string) error
	MarkEmailVerified(ctx context.Context, userID int64, at time.Time) error
}

type userProvider struct{ db postgres.IBaseProvider }

func NewUserProvider(db postgres.IBaseProvider) IUserProvider {
	return &userProvider{db: db}
}

func (p *userProvider) Create(ctx context.Context, u entity.User) (*entity.User, error) {
	row := models.User{
		PublicID:        u.PublicID,
		Email:           u.Email,
		PasswordHash:    u.PasswordHash,
		EmailVerifiedAt: u.EmailVerifiedAt,
	}
	if _, err := p.db.Conn(ctx).NewInsert().
		Model(&row).
		Returning("id, public_id, created_at, updated_at").
		Exec(ctx); err != nil {
		return nil, err
	}
	e := userToEntity(row)
	return &e, nil
}

// GetByEmail сравнивает по LOWER(email) — таблица уникальна по тому же
// выражению (см. 0005_users.up.sql).
func (p *userProvider) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := models.User{}
	err := p.db.Conn(ctx).NewSelect().
		Model(&row).
		Where("LOWER(email) = LOWER(?)", email).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := userToEntity(row)
	return &e, nil
}

func (p *userProvider) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	row := models.User{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := userToEntity(row)
	return &e, nil
}

func (p *userProvider) GetByPublicID(ctx context.Context, publicID string) (*entity.User, error) {
	row := models.User{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("public_id = ?", publicID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e := userToEntity(row)
	return &e, nil
}

func (p *userProvider) UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string) error {
	_, err := p.db.Conn(ctx).NewUpdate().
		Model((*models.User)(nil)).
		Set("password_hash = ?", passwordHash).
		Set("updated_at = NOW()").
		Where("id = ?", userID).
		Exec(ctx)
	return err
}

func (p *userProvider) MarkEmailVerified(ctx context.Context, userID int64, at time.Time) error {
	_, err := p.db.Conn(ctx).NewUpdate().
		Model((*models.User)(nil)).
		Set("email_verified_at = ?", at).
		Set("updated_at = NOW()").
		Where("id = ?", userID).
		Exec(ctx)
	return err
}