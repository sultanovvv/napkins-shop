package product

import (
	"context"
	"errors"
	"sort"

	"go.uber.org/zap"

	attributeuc "napkins-shop/public-executor/internal/usecases/attribute"
	categoryuc "napkins-shop/public-executor/internal/usecases/category"
	"shared/entity"
	core_db "shared/providers/core-db"
)

var ErrNotFound = errors.New("product not found")

type IProductUseCases interface {
	GetProductsList(ctx context.Context, in GetProductsListInUDTO) ([]entity.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*entity.Product, error)
	GetProductByID(ctx context.Context, id int64) (*entity.Product, error)
}

type useCases struct {
	logger      *zap.Logger
	productProv core_db.IProductProvider
	attrValProv core_db.IProductAttributeProvider
	imageProv   core_db.IProductImageProvider
	categoryUC  categoryuc.ICategoryUseCases
	attributeUC attributeuc.IAttributeUseCases
}

func NewUseCase(
	productProv core_db.IProductProvider,
	attrValProv core_db.IProductAttributeProvider,
	imageProv core_db.IProductImageProvider,
	categoryUC categoryuc.ICategoryUseCases,
	attributeUC attributeuc.IAttributeUseCases,
	logger *zap.Logger,
) IProductUseCases {
	return &useCases{
		productProv: productProv,
		attrValProv: attrValProv,
		imageProv:   imageProv,
		categoryUC:  categoryUC,
		attributeUC: attributeUC,
		logger:      logger,
	}
}

func (u *useCases) GetProductsList(ctx context.Context, in GetProductsListInUDTO) ([]entity.Product, error) {
	rows, err := u.productProv.List(ctx, core_db.ProductListFilter{CategorySlug: in.CategorySlug})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []entity.Product{}, nil
	}

	catByID, err := u.categoriesByID(ctx)
	if err != nil {
		return nil, err
	}
	attrByID, err := u.attributesByID(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	values, err := u.attrValProv.ListByProducts(ctx, ids)
	if err != nil {
		return nil, err
	}
	valuesByProduct := groupAttributes(values, attrByID)

	images, err := u.imageProv.ListByProducts(ctx, ids)
	if err != nil {
		return nil, err
	}
	imagesByProduct := groupImages(images)

	out := make([]entity.Product, len(rows))
	for i := range rows {
		out[i] = enrich(rows[i], catByID, valuesByProduct[rows[i].ID], imagesByProduct[rows[i].ID])
	}
	return out, nil
}

func (u *useCases) GetProductBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	row, err := u.productProv.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}
	return u.hydrate(ctx, *row)
}

func (u *useCases) GetProductByID(ctx context.Context, id int64) (*entity.Product, error) {
	row, err := u.productProv.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}
	return u.hydrate(ctx, *row)
}

func (u *useCases) hydrate(ctx context.Context, p entity.Product) (*entity.Product, error) {
	catByID, err := u.categoriesByID(ctx)
	if err != nil {
		return nil, err
	}
	attrByID, err := u.attributesByID(ctx)
	if err != nil {
		return nil, err
	}
	values, err := u.attrValProv.ListByProduct(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	images, err := u.imageProv.ListByProduct(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	out := enrich(p, catByID, resolveAttributes(values, attrByID), images)
	return &out, nil
}

func (u *useCases) categoriesByID(ctx context.Context) (map[int64]*entity.Category, error) {
	cats, err := u.categoryUC.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[int64]*entity.Category, len(cats))
	for i := range cats {
		c := cats[i]
		out[c.ID] = &c
	}
	return out, nil
}

func (u *useCases) attributesByID(ctx context.Context) (map[int64]entity.Attribute, error) {
	attrs, err := u.attributeUC.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[int64]entity.Attribute, len(attrs))
	for _, a := range attrs {
		out[a.ID] = a
	}
	return out, nil
}

func enrich(
	p entity.Product,
	catByID map[int64]*entity.Category,
	attrs []entity.AttributeValue,
	images []entity.Image,
) entity.Product {
	p.Attributes = attrs
	p.Images = images
	if p.CategoryID != nil {
		if c, ok := catByID[*p.CategoryID]; ok {
			p.Category = c
		}
	}
	return p
}

// groupAttributes раскладывает плоский список значений по продукту с
// сохранением SortOrder каталога атрибутов.
func groupAttributes(
	rows []entity.AttributeValueRow,
	attrByID map[int64]entity.Attribute,
) map[int64][]entity.AttributeValue {
	grouped := make(map[int64][]entity.AttributeValue, len(rows))
	for _, r := range rows {
		attr, ok := attrByID[r.AttributeID]
		if !ok {
			continue
		}
		grouped[r.ProductID] = append(grouped[r.ProductID], entity.AttributeValue{
			Attribute: attr,
			ValueText: r.ValueText,
			ValueInt:  r.ValueInt,
		})
	}
	for k := range grouped {
		sort.SliceStable(grouped[k], func(i, j int) bool {
			return grouped[k][i].Attribute.SortOrder < grouped[k][j].Attribute.SortOrder
		})
	}
	return grouped
}

func resolveAttributes(
	rows []entity.AttributeValueRow,
	attrByID map[int64]entity.Attribute,
) []entity.AttributeValue {
	out := make([]entity.AttributeValue, 0, len(rows))
	for _, r := range rows {
		attr, ok := attrByID[r.AttributeID]
		if !ok {
			continue
		}
		out = append(out, entity.AttributeValue{
			Attribute: attr,
			ValueText: r.ValueText,
			ValueInt:  r.ValueInt,
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Attribute.SortOrder < out[j].Attribute.SortOrder
	})
	return out
}

func groupImages(images []entity.Image) map[int64][]entity.Image {
	out := make(map[int64][]entity.Image, len(images))
	for _, img := range images {
		out[img.ProductID] = append(out[img.ProductID], img)
	}
	return out
}