package attribute

import (
	"errors"

	"go.uber.org/zap"

	attributerepo "napkins-shop/public-executor/internal/repository/attribute"
	"shared/entity"
)

var ErrNotFound = errors.New("attribute not found")

type IAttributeUseCases interface {
	List() ([]entity.Attribute, error)
	GetByID(id int64) (*entity.Attribute, error)
	GetBySlug(slug string) (*entity.Attribute, error)
}

type useCases struct {
	logger *zap.Logger
	repo   attributerepo.IAttributeRepository
}

func NewUseCase(repo attributerepo.IAttributeRepository, logger *zap.Logger) IAttributeUseCases {
	return &useCases{repo: repo, logger: logger}
}

func (u *useCases) List() ([]entity.Attribute, error) {
	return u.repo.List()
}

func (u *useCases) GetByID(id int64) (*entity.Attribute, error) {
	a, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	return a, nil
}

func (u *useCases) GetBySlug(slug string) (*entity.Attribute, error) {
	a, err := u.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	return a, nil
}
