package product_image

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	productrepo "napkins-shop/public-executor/internal/repository/product"
	productimagerepo "napkins-shop/public-executor/internal/repository/product_image"
	"napkins-shop/public-executor/internal/storage/s3gateway"
	"napkins-shop/public-executor/internal/storage/s3url"
	"shared/entity"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidImage    = errors.New("invalid image")
)

// UploadResult — ответ usecase'у; собирается из созданной строки + URL,
// который presenter уже знает как собрать.
type UploadResult struct {
	ID        int64
	Key       string
	URL       string
	IsPrimary bool
	SortOrder int
}

type IUploader interface {
	Upload(ctx context.Context, in UploadInUDTO) (*UploadResult, error)
}

type uploader struct {
	logger      *zap.Logger
	productRepo productrepo.IProductRepository
	imageRepo   productimagerepo.IProductImageRepository
	gateway     s3gateway.Gateway
	urls        s3url.Builder
}

func NewUploader(
	productRepo productrepo.IProductRepository,
	imageRepo productimagerepo.IProductImageRepository,
	gateway s3gateway.Gateway,
	urls s3url.Builder,
	logger *zap.Logger,
) IUploader {
	return &uploader{
		logger:      logger,
		productRepo: productRepo,
		imageRepo:   imageRepo,
		gateway:     gateway,
		urls:        urls,
	}
}

func (u *uploader) Upload(ctx context.Context, in UploadInUDTO) (*UploadResult, error) {
	if in.ProductSlug == "" || in.Body == nil || in.Size <= 0 {
		return nil, ErrInvalidImage
	}

	product, err := u.productRepo.GetBySlug(in.ProductSlug)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	key := buildKey(product.Slug, in.Filename)

	if err := u.gateway.Put(ctx, key, in.Body, in.Size, in.ContentType); err != nil {
		return nil, fmt.Errorf("s3 put: %w", err)
	}

	if in.IsPrimary {
		if err := u.imageRepo.ClearPrimary(product.ID); err != nil {
			// Откатим только что залитый объект, чтобы не оставлять
			// «висячую» картинку в бакете без записи в БД.
			_ = u.gateway.Delete(ctx, key)
			return nil, fmt.Errorf("clear primary: %w", err)
		}
	}

	img := &entity.Image{
		ProductID: product.ID,
		Key:       key,
		SortOrder: in.SortOrder,
		IsPrimary: in.IsPrimary,
	}
	if err := u.imageRepo.Insert(img); err != nil {
		_ = u.gateway.Delete(ctx, key)
		return nil, fmt.Errorf("insert image row: %w", err)
	}

	return &UploadResult{
		ID:        img.ID,
		Key:       img.Key,
		URL:       u.urls.URL(img.Key),
		IsPrimary: img.IsPrimary,
		SortOrder: img.SortOrder,
	}, nil
}

// buildKey собирает stable-ключ products/{slug}/{timestamp-suffix}.{ext}.
// Уникальный суффикс защищает от коллизий имён при множественных загрузках.
func buildKey(slug, filename string) string {
	ext := strings.ToLower(extOf(filename))
	if ext == "" {
		ext = ".bin"
	}
	ts := time.Now().UTC().Format("20060102-150405")
	return fmt.Sprintf("products/%s/%s%s", slug, ts, ext)
}

func extOf(filename string) string {
	i := strings.LastIndex(filename, ".")
	if i < 0 {
		return ""
	}
	return filename[i:]
}
