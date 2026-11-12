package core_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"shared/configs/postgres"
	"shared/providers/core-db/models"
)

func (p *backofficeCategoryProvider) write(ctx context.Context, callback postgres.TransactionCallback) error {
	err := p.tx.Run(ctx, func(ctx context.Context) error {
		// Для небольшого справочника сериализуем изменения дерева одной
		// табличной блокировкой. Две одновременные смены родителей не могут
		// пройти проверку по старому дереву и совместно создать цикл.
		// Обычные SELECT продолжают выполняться. Используется READ COMMITTED
		// (стандартный уровень TransactionProvider); блокировка — до чтений.
		if err := p.lockTree(ctx); err != nil {
			return err
		}
		return callback(ctx)
	})
	return backofficeCategoryError(err)
}

func (p *backofficeCategoryProvider) lockTree(ctx context.Context) error {
	_, err := p.db.Conn(ctx).NewRaw("LOCK TABLE categories IN SHARE ROW EXCLUSIVE MODE").Exec(ctx)
	return err
}

func (p *backofficeCategoryProvider) getForUpdate(ctx context.Context, id int64) (*models.Category, error) {
	row := models.Category{}
	err := p.db.Conn(ctx).NewSelect().Model(&row).Where("c.id = ?", id).For("UPDATE").Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCategoryNotFound
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *backofficeCategoryProvider) validateParent(ctx context.Context, categoryID int64, parentID *int64) error {
	if parentID == nil {
		return nil
	}
	if *parentID == categoryID {
		return ErrCategoryCycle
	}
	parent, err := p.GetByID(ctx, *parentID)
	if err != nil {
		return err
	}
	if parent == nil {
		return ErrCategoryParentNotFound
	}

	var invalid bool
	err = p.db.Conn(ctx).NewRaw(`
		WITH RECURSIVE ancestors AS (
			SELECT id, parent_id FROM categories WHERE id = ?
			UNION
			SELECT c.id, c.parent_id
			FROM categories AS c
			JOIN ancestors AS a ON c.id = a.parent_id
		)
		SELECT EXISTS (SELECT 1 FROM ancestors WHERE id = ?)
			OR NOT EXISTS (SELECT 1 FROM ancestors WHERE parent_id IS NULL)
	`, *parentID, categoryID).Scan(ctx, &invalid)
	if err != nil {
		return err
	}
	// UNION завершает обход даже при уже испорченном циклом дереве.
	// У корректной цепочки предков обязательно есть корень с parent_id NULL.
	if invalid {
		return ErrCategoryCycle
	}
	return nil
}

func backofficeCategoryError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch {
		case pgErr.Code == "23505" && pgErr.ConstraintName == "categories_slug_key":
			return fmt.Errorf("%w: %w", ErrCategorySlugConflict, err)
		case pgErr.Code == "23503" && pgErr.ConstraintName == "categories_parent_id_fkey":
			return fmt.Errorf("%w: %w", ErrCategoryParentNotFound, err)
		}
	}
	return err
}
