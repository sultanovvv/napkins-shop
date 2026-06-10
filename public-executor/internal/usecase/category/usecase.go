package category

import (
	"errors"

	"go.uber.org/zap"

	categoryrepo "napkins-shop/public-executor/internal/repository/category"
	"shared/entity"
)

var ErrNotFound = errors.New("category not found")

type ICategoryUseCases interface {
	GetTree() ([]entity.Category, error)
	GetBySlug(slug string) (*entity.Category, error)
	GetByID(id int64) (*entity.Category, error)
	List() ([]entity.Category, error)
}

type useCases struct {
	logger *zap.Logger
	repo   categoryrepo.ICategoryRepository
}

func NewUseCase(repo categoryrepo.ICategoryRepository, logger *zap.Logger) ICategoryUseCases {
	return &useCases{repo: repo, logger: logger}
}

func (u *useCases) GetTree() ([]entity.Category, error) {
	rows, err := u.repo.List()
	if err != nil {
		return nil, err
	}

	byID := make(map[int64]*entity.Category, len(rows))
	for i := range rows {
		// Копию делаем, чтобы byID указывал на новый аллок, а не на
		// элемент локального слайса (тот можно затем переиспользовать).
		c := rows[i]
		byID[c.ID] = &c
	}

	var roots []*entity.Category
	for i := range rows {
		node := byID[rows[i].ID]
		if rows[i].ParentID == nil {
			roots = append(roots, node)
			continue
		}
		parent, ok := byID[*rows[i].ParentID]
		if !ok {
			// orphaned child — surface it as root so it's not lost
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, *node)
	}

	out := make([]entity.Category, len(roots))
	for i, r := range roots {
		out[i] = *r
	}
	return out, nil
}

func (u *useCases) GetBySlug(slug string) (*entity.Category, error) {
	c, err := u.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	return c, nil
}

func (u *useCases) GetByID(id int64) (*entity.Category, error) {
	c, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	return c, nil
}

func (u *useCases) List() ([]entity.Category, error) {
	return u.repo.List()
}
