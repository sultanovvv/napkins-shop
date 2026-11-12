// Package slug содержит общие правила slug для административных сценариев.
package slug

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"shared/utils"
)

var (
	ErrInvalid   = errors.New("invalid slug")
	ErrExhausted = errors.New("slug suffix range exhausted")
	validSlug    = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
)

// Normalize проверяет ручной slug, не выполняя транслитерацию.
func Normalize(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if !validSlug.MatchString(value) {
		return "", ErrInvalid
	}
	return value, nil
}

// Generate подбирает первый свободный slug: base, base-2, base-3, ... .
// isTaken определяет область уникальности (категории, товары, атрибуты).
// Результат не резервируется: вызывающий сценарий обрабатывает конфликт INSERT.
func Generate(ctx context.Context, name, fallback string, isTaken func(context.Context, string) (bool, error)) (string, error) {
	base := utils.Slug(name)
	if base == "" {
		base = utils.Slug(fallback)
	}
	if base == "" {
		return "", ErrInvalid
	}
	for suffix := int64(1); suffix > 0; suffix++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		candidate := base
		if suffix > 1 {
			candidate += "-" + strconv.FormatInt(suffix, 10)
		}
		taken, err := isTaken(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !taken {
			return candidate, nil
		}
	}
	return "", ErrExhausted
}
