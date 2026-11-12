package category

import (
	"strings"
	"unicode/utf8"

	sluguc "napkins-shop/backoffice-executor/internal/usecases/slug"
)

func normalizeName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", validationError("name", "must not be empty")
	}
	if !utf8.ValidString(value) || strings.ContainsRune(value, 0) {
		return "", validationError("name", "must be valid text without null characters")
	}
	return value, nil
}

func normalizeSlug(value string) (string, error) {
	value, err := sluguc.Normalize(value)
	if err != nil {
		return "", validationError("slug", "must contain lowercase Latin letters, digits and single separating hyphens")
	}
	return value, nil
}

func validateID(field string, id int64) error {
	if id <= 0 {
		return validationError(field, "must be positive")
	}
	return nil
}

func validateParentID(id *int64) error {
	if id != nil {
		return validateID("parentId", *id)
	}
	return nil
}

func validateSortOrder(value int) error {
	if value < -1<<31 || value > 1<<31-1 {
		return validationError("sortOrder", "must fit a signed 32-bit integer")
	}
	return nil
}

func normalizeUpdate(in UpdateCategoryInUDTO) (UpdateCategoryInUDTO, error) {
	if err := validateID("id", in.ID); err != nil {
		return in, err
	}
	if in.Name.Set {
		if in.Name.Value == nil {
			return in, validationError("name", "must not be null")
		}
		name, err := normalizeName(*in.Name.Value)
		if err != nil {
			return in, err
		}
		in.Name.Value = &name
	}
	if in.Slug.Set {
		if in.Slug.Value == nil {
			return in, validationError("slug", "must not be null")
		}
		slug, err := normalizeSlug(*in.Slug.Value)
		if err != nil {
			return in, err
		}
		in.Slug.Value = &slug
	}
	if in.ClearParentID && in.ParentID != nil {
		return in, validationError("clearParentId", "cannot be used together with parentId")
	}
	if err := validateParentID(in.ParentID); err != nil {
		return in, err
	}
	if in.SortOrder.Set {
		if in.SortOrder.Value == nil {
			return in, validationError("sortOrder", "must not be null")
		}
		if err := validateSortOrder(*in.SortOrder.Value); err != nil {
			return in, err
		}
	}
	return in, nil
}
