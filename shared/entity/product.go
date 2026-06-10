package entity

import "time"

// Product — основная сущность каталога.
//
// Поля CategoryID/Category, Attributes, Images заполняются разными слоями:
// репозиторий кладёт raw-данные (CategoryID), usecase подмешивает
// связанные сущности (Category, Attributes, Images). Поэтому пустой слайс
// или nil в составных полях — это "ещё не собрали", а не "нет данных".
type Product struct {
	ID         int64
	Slug       string
	Name       string
	CategoryID *int64
	Category   *Category
	Attributes []AttributeValue
	Images     []Image
}

// Image — запись о картинке товара. Key — путь объекта в S3-совместимом
// хранилище; URL строит presenter (см. internal/storage/s3url).
type Image struct {
	ID        int64
	ProductID int64
	Key       string
	SortOrder int
	IsPrimary bool
	CreatedAt time.Time
}
