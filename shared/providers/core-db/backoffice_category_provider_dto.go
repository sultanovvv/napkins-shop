package core_db

// BackofficeCreateCategoryInPDTO содержит готовые к сохранению поля категории.
// ParentID == nil означает корневую категорию. Уровень отдельно не хранится.
type BackofficeCreateCategoryInPDTO struct {
	Name      string
	Slug      string
	ParentID  *int64
	SortOrder int
}

// BackofficeUpdateCategoryInPDTO описывает итоговое состояние категории,
// включая неизменённые поля. Нулевой SortOrder и nil ParentID записываются явно.
type BackofficeUpdateCategoryInPDTO struct {
	ID        int64
	Name      string
	Slug      string
	ParentID  *int64
	SortOrder int
}
