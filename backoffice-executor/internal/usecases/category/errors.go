package category

import (
	"errors"

	"napkins-shop/backoffice-executor/internal/customerr"
	core_db "shared/providers/core-db"
)

func categoryNotFound(cause error) *customerr.Error {
	return fieldError(customerr.NewNotFound("Category not found", cause), "CATEGORY_NOT_FOUND", "id")
}

func validationError(field, message string) *customerr.Error {
	err := customerr.NewValidation("Request validation failed", nil)
	err.Details = []customerr.Detail{{Code: "INVALID_VALUE", Message: message, Target: field}}
	return err
}

func fieldError(err *customerr.Error, code, target string) *customerr.Error {
	err.Details = []customerr.Detail{{Code: code, Message: err.Message, Target: target}}
	return err
}

// categoryError формирует прикладную ошибку на границе с провайдером БД.
func categoryError(err error) error {
	if err == nil {
		return nil
	}
	var custom *customerr.Error
	if errors.As(err, &custom) {
		return err
	}
	switch {
	case errors.Is(err, core_db.ErrCategoryNotFound):
		return categoryNotFound(err)
	case errors.Is(err, core_db.ErrCategoryParentNotFound):
		return fieldError(customerr.NewValidation("Parent category not found", err), "CATEGORY_PARENT_NOT_FOUND", "parentId")
	case errors.Is(err, core_db.ErrCategorySlugConflict):
		return fieldError(customerr.NewConflict("Slug already exists", err), "SLUG_ALREADY_EXISTS", "slug")
	case errors.Is(err, core_db.ErrCategoryCycle):
		return fieldError(customerr.NewValidation("Parent would create a category hierarchy cycle", err), "CATEGORY_CYCLE", "parentId")
	default:
		return customerr.NewInternal("Internal Server Error", err)
	}
}
