package category

type CreateCategoryInUDTO struct {
	Name      string
	Slug      *string
	ParentID  *int64
	SortOrder int
}

// PatchField различает отсутствующее поле, явный null и заданное значение.
// HTTP-слой должен сохранять это различие при преобразовании входного запроса.
type PatchField[T any] struct {
	Set   bool
	Value *T
}

type UpdateCategoryInUDTO struct {
	ID            int64
	Name          PatchField[string]
	Slug          PatchField[string]
	ParentID      *int64
	ClearParentID bool
	SortOrder     PatchField[int]
}

// CategoryListItemOutUDTO — строка списка в порядке обхода дерева.
// Структура задаётся ParentID, отдельного поля уровня нет.
type CategoryListItemOutUDTO struct {
	ID        int64
	Name      string
	Slug      string
	ParentID  *int64
	SortOrder int
}

type CategoryNodeOutUDTO struct {
	CategoryListItemOutUDTO
	Children []*CategoryNodeOutUDTO
}
