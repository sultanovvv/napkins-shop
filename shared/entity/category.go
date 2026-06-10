package entity

// Category — узел каталога. Children заполняется только при построении дерева
// (см. usecase GetTree). На листьях Children — пустой слайс.
type Category struct {
	ID        int64
	Slug      string
	Name      string
	ParentID  *int64
	SortOrder int
	Children  []Category
}
