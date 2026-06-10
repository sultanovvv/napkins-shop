package product

// GetProductsListInUDTO — параметры запроса списка товаров.
// CategorySlug пустой = вся витрина; иначе фильтруется по slug категории
// (включая всех её потомков, если slug указывает на родительский узел).
type GetProductsListInUDTO struct {
	CategorySlug string
}
