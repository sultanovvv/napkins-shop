package core_db

import (
	"context"
	"errors"

	"shared/configs/postgres"
	"shared/entity"
	"shared/providers/core-db/models"
)

var (
	ErrCategoryNotFound       = errors.New("category not found")
	ErrCategoryParentNotFound = errors.New("category parent not found")
	ErrCategorySlugConflict   = errors.New("category slug already exists")
	ErrCategoryCycle          = errors.New("category hierarchy cycle")
)

// IBackofficeCategoryProvider расширяет чтение категорий операциями записи.
// Slug передаётся готовым: генерация и валидация входных полей относятся к usecase.
type IBackofficeCategoryProvider interface {
	ICategoryProvider
	// GetByIDForUpdate вызывается внутри ITransactionProvider.Run перед
	// сборкой частичного обновления. Сначала блокирует изменения дерева.
	GetByIDForUpdate(ctx context.Context, id int64) (*entity.Category, error)
	Create(ctx context.Context, in BackofficeCreateCategoryInPDTO) (*entity.Category, error)
	Update(ctx context.Context, in BackofficeUpdateCategoryInPDTO) (*entity.Category, error)
	Delete(ctx context.Context, id int64) error
}

func (p *backofficeCategoryProvider) GetByIDForUpdate(ctx context.Context, id int64) (*entity.Category, error) {
	if err := p.lockTree(ctx); err != nil {
		return nil, err
	}
	row, err := p.getForUpdate(ctx, id)
	if err != nil {
		return nil, err
	}
	category := categoryToEntity(*row)
	return &category, nil
}

type backofficeCategoryProvider struct {
	categoryProvider
	tx ITransactionProvider
}

var _ IBackofficeCategoryProvider = (*backofficeCategoryProvider)(nil)

// NewBackofficeCategoryProvider использует BaseProvider и TransactionProvider
// одного соединения core-db. GetByID/GetBySlug переиспользуют общий провайдер.
func NewBackofficeCategoryProvider(db postgres.IBaseProvider, tx ITransactionProvider) IBackofficeCategoryProvider {
	return &backofficeCategoryProvider{
		categoryProvider: categoryProvider{db: db},
		tx:               tx,
	}
}

// List возвращает плоский список. Вложенность задаётся ParentID;
// построение дерева выполняется на прикладном слое.
func (p *backofficeCategoryProvider) List(ctx context.Context) ([]entity.Category, error) {
	var rows []models.Category
	if err := p.db.Conn(ctx).NewSelect().Model(&rows).
		OrderExpr("parent_id NULLS FIRST").
		Order("sort_order", "name", "id").
		Scan(ctx); err != nil {
		return nil, err
	}

	out := make([]entity.Category, len(rows))
	for i := range rows {
		out[i] = categoryToEntity(rows[i])
	}
	return out, nil
}

func (p *backofficeCategoryProvider) Create(ctx context.Context, in BackofficeCreateCategoryInPDTO) (*entity.Category, error) {
	var created entity.Category
	err := p.write(ctx, func(ctx context.Context) error {
		if err := p.validateParent(ctx, 0, in.ParentID); err != nil {
			return err
		}

		row := models.Category{
			Name: in.Name, Slug: in.Slug, ParentID: in.ParentID, SortOrder: in.SortOrder,
		}
		if _, err := p.db.Conn(ctx).NewInsert().Model(&row).Returning("*").Exec(ctx); err != nil {
			return err
		}
		created = categoryToEntity(row)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &created, nil
}

// Update полностью заменяет редактируемые поля. ParentID == nil переносит
// категорию в корень; для PATCH usecase сначала собирает итоговое состояние.
func (p *backofficeCategoryProvider) Update(ctx context.Context, in BackofficeUpdateCategoryInPDTO) (*entity.Category, error) {
	var updated entity.Category
	err := p.write(ctx, func(ctx context.Context) error {
		row, err := p.getForUpdate(ctx, in.ID)
		if err != nil {
			return err
		}
		if err := p.validateParent(ctx, in.ID, in.ParentID); err != nil {
			return err
		}

		row.Name = in.Name
		row.Slug = in.Slug
		row.ParentID = in.ParentID
		row.SortOrder = in.SortOrder
		if _, err := p.db.Conn(ctx).NewUpdate().Model(row).
			Column("name", "slug", "parent_id", "sort_order").
			WherePK().Returning("*").Exec(ctx); err != nil {
			return err
		}
		updated = categoryToEntity(*row)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete удаляет категорию со всем поддеревом через ON DELETE CASCADE.
// ON DELETE SET NULL сохраняет товары всех удалённых узлов без категории.
func (p *backofficeCategoryProvider) Delete(ctx context.Context, id int64) error {
	return p.write(ctx, func(ctx context.Context) error {
		if _, err := p.getForUpdate(ctx, id); err != nil {
			return err
		}

		_, err := p.db.Conn(ctx).NewDelete().Model((*models.Category)(nil)).Where("id = ?", id).Exec(ctx)
		return err
	})
}
