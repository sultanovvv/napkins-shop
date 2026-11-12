package category

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"napkins-shop/backoffice-executor/internal/customerr"
	sluguc "napkins-shop/backoffice-executor/internal/usecases/slug"
	"shared/entity"
	core_db "shared/providers/core-db"
)

const maxSlugCreateAttempts = 5

type ICategoryUseCases interface {
	Create(ctx context.Context, in CreateCategoryInUDTO) (*entity.Category, error)
	Update(ctx context.Context, in UpdateCategoryInUDTO) (*entity.Category, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	List(ctx context.Context) ([]CategoryListItemOutUDTO, error)
	GetTree(ctx context.Context) ([]*CategoryNodeOutUDTO, error)
}

type useCases struct {
	provider core_db.IBackofficeCategoryProvider
	tx       core_db.ITransactionProvider
}

var _ ICategoryUseCases = (*useCases)(nil)

func NewUseCase(provider core_db.IBackofficeCategoryProvider, tx core_db.ITransactionProvider) ICategoryUseCases {
	return &useCases{provider: provider, tx: tx}
}

func (u *useCases) Create(ctx context.Context, in CreateCategoryInUDTO) (*entity.Category, error) {
	name, err := normalizeName(in.Name)
	if err != nil {
		return nil, err
	}
	if err := validateParentID(in.ParentID); err != nil {
		return nil, err
	}
	if err := validateSortOrder(in.SortOrder); err != nil {
		return nil, err
	}
	fields := core_db.BackofficeCreateCategoryInPDTO{
		Name: name, ParentID: in.ParentID, SortOrder: in.SortOrder,
	}
	if in.Slug != nil && strings.TrimSpace(*in.Slug) != "" {
		fields.Slug, err = normalizeSlug(*in.Slug)
		if err != nil {
			return nil, err
		}
		created, err := u.provider.Create(ctx, fields)
		return created, categoryError(err)
	}

	var lastConflict error
	for attempt := 0; attempt < maxSlugCreateAttempts; attempt++ {
		fields.Slug, err = sluguc.Generate(ctx, name, "category", u.slugTaken)
		if err != nil {
			return nil, err
		}
		created, err := u.provider.Create(ctx, fields)
		if err == nil {
			return created, nil
		}
		if !errors.Is(err, core_db.ErrCategorySlugConflict) {
			return nil, categoryError(err)
		}
		// Провайдер уже откатил транзакцию/сохранённую точку. Если между
		// подбором и INSERT slug заняли, повторяем с новым свободным значением.
		lastConflict = err
	}
	return nil, customerr.NewUnavailable("Service Unavailable", fmt.Errorf("could not allocate category slug: %w", lastConflict))
}

func (u *useCases) Update(ctx context.Context, in UpdateCategoryInUDTO) (*entity.Category, error) {
	in, err := normalizeUpdate(in)
	if err != nil {
		return nil, err
	}
	var updated *entity.Category
	err = u.tx.Run(ctx, func(ctx context.Context) error {
		// Читаем после блокировки дерева, до объединения PATCH-полей:
		// параллельное изменение не затрётся устаревшим полным DTO.
		current, err := u.provider.GetByIDForUpdate(ctx, in.ID)
		if err != nil {
			return err
		}
		if !in.Name.Set && !in.Slug.Set && in.ParentID == nil && !in.ClearParentID && !in.SortOrder.Set {
			updated = current
			return nil
		}
		fields := core_db.BackofficeUpdateCategoryInPDTO{
			ID: current.ID, Name: current.Name, Slug: current.Slug,
			ParentID: current.ParentID, SortOrder: current.SortOrder,
		}
		if in.Name.Set {
			fields.Name = *in.Name.Value
		}
		if in.Slug.Set {
			fields.Slug = *in.Slug.Value
		}
		if in.ClearParentID {
			fields.ParentID = nil
		} else if in.ParentID != nil {
			fields.ParentID = in.ParentID
		}
		if in.SortOrder.Set {
			fields.SortOrder = *in.SortOrder.Value
		}
		updated, err = u.provider.Update(ctx, fields)
		return err
	})
	if err != nil {
		return nil, categoryError(err)
	}
	return updated, nil
}

func (u *useCases) Delete(ctx context.Context, id int64) error {
	if err := validateID("id", id); err != nil {
		return err
	}
	return categoryError(u.provider.Delete(ctx, id))
}

func (u *useCases) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	if err := validateID("id", id); err != nil {
		return nil, err
	}
	category, err := u.provider.GetByID(ctx, id)
	if err != nil {
		return nil, categoryError(err)
	}
	if category == nil {
		return nil, categoryNotFound(nil)
	}
	return category, nil
}

func (u *useCases) List(ctx context.Context) ([]CategoryListItemOutUDTO, error) {
	rows, err := u.provider.List(ctx)
	if err != nil {
		return nil, categoryError(err)
	}
	return orderedCategories(rows)
}

func (u *useCases) GetTree(ctx context.Context) ([]*CategoryNodeOutUDTO, error) {
	items, err := u.List(ctx)
	if err != nil {
		return nil, err
	}
	roots := make([]*CategoryNodeOutUDTO, 0)
	byID := make(map[int64]*CategoryNodeOutUDTO, len(items))
	for _, item := range items {
		node := &CategoryNodeOutUDTO{
			CategoryListItemOutUDTO: item,
			Children:                make([]*CategoryNodeOutUDTO, 0),
		}
		byID[item.ID] = node
		if item.ParentID == nil {
			roots = append(roots, node)
		} else {
			// orderedCategories гарантирует, что родитель уже добавлен.
			parent := byID[*item.ParentID]
			parent.Children = append(parent.Children, node)
		}
	}
	return roots, nil
}

func (u *useCases) slugTaken(ctx context.Context, value string) (bool, error) {
	category, err := u.provider.GetBySlug(ctx, value)
	return category != nil, categoryError(err)
}
