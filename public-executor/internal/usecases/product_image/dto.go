package product_image

import "io"

// UploadInUDTO — параметры заливки картинки товара.
// Body читается потоково; usecase не закрывает Reader — это ответственность
// вызывающего (хендлер).
type UploadInUDTO struct {
	ProductSlug string
	Filename    string
	ContentType string
	Body        io.Reader
	Size        int64
	IsPrimary   bool
	SortOrder   int
}
