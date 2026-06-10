package product

import (
	"errors"
	"sort"

	"go.uber.org/zap"

	productrepo "napkins-shop/public-executor/internal/repository/product"
	productattrrepo "napkins-shop/public-executor/internal/repository/product_attribute"
	productimagerepo "napkins-shop/public-executor/internal/repository/product_image"
	attributeuc "napkins-shop/public-executor/internal/usecase/attribute"
	categoryuc "napkins-shop/public-executor/internal/usecase/category"
	"shared/entity"
)

var ErrNotFound = errors.New("product not found")

type IProductUseCases interface {
	GetProductsList(in GetProductsListInUDTO) ([]entity.Product, error)
	GetProductBySlug(slug string) (*entity.Product, error)
	GetProductByID(id int64) (*entity.Product, error)
}

type useCases struct {
	logger        *zap.Logger
	productRepo   productrepo.IProductRepository
	attrValueRepo productattrrepo.IProductAttributeRepository
	imageRepo     productimagerepo.IProductImageRepository
	categoryUC    categoryuc.ICategoryUseCases
	attributeUC   attributeuc.IAttributeUseCases
}

func NewUseCase(
	productRepo productrepo.IProductRepository,
	attrValueRepo productattrrepo.IProductAttributeRepository,
	imageRepo productimagerepo.IProductImageRepository,
	categoryUC categoryuc.ICategoryUseCases,
	attributeUC attributeuc.IAttributeUseCases,
	logger *zap.Logger,
) IProductUseCases {
	return &useCases{
		productRepo:   productRepo,
		attrValueRepo: attrValueRepo,
		imageRepo:     imageRepo,
		categoryUC:    categoryUC,
		attributeUC:   attributeUC,
		logger:        logger,
	}
}

func (u *useCases) GetProductsList(in GetProductsListInUDTO) ([]entity.Product, error) {
	rows, err := u.productRepo.List(productrepo.ListFilter{CategorySlug: in.CategorySlug})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []entity.Product{}, nil
	}

	catByID, err := u.categoriesByID()
	if err != nil {
		return nil, err
	}
	attrByID, err := u.attributesByID()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	values, err := u.attrValueRepo.ListByProducts(ids)
	if err != nil {
		return nil, err
	}
	valuesByProduct := groupAttributes(values, attrByID)

	images, err := u.imageRepo.ListByProducts(ids)
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

func (u *useCases) GetProductBySlug(slug string) (*entity.Product, error) {
	row, err := u.productRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}
	return u.hydrate(*row)
}

func (u *useCases) GetProductByID(id int64) (*entity.Product, error) {
	row, err := u.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}
	return u.hydrate(*row)
}

func (u *useCases) hydrate(p entity.Product) (*entity.Product, error) {
	catByID, err := u.categoriesByID()
	if err != nil {
		return nil, err
	}
	attrByID, err := u.attributesByID()
	if err != nil {
		return nil, err
	}
	values, err := u.attrValueRepo.ListByProduct(p.ID)
	if err != nil {
		return nil, err
	}
	images, err := u.imageRepo.ListByProduct(p.ID)
	if err != nil {
		return nil, err
	}
	out := enrich(p, catByID, resolveAttributes(values, attrByID), images)
	return &out, nil
}

func (u *useCases) categoriesByID() (map[int64]*entity.Category, error) {
	cats, err := u.categoryUC.List()
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

func (u *useCases) attributesByID() (map[int64]entity.Attribute, error) {
	attrs, err := u.attributeUC.List()
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
