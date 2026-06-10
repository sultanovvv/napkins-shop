package entity

// ValueType — закрытый список типов значения атрибута. Контракт между
// каталогом атрибутов и хранилищем значений: ровно одно из value_text /
// value_int должно быть заполнено в product_attribute_values.
type ValueType string

const (
	ValueTypeString ValueType = "string"
	ValueTypeInt    ValueType = "int"
)

// Attribute — справочник атрибута (size/layers/color/...).
type Attribute struct {
	ID        int64
	Slug      string
	Name      string
	ValueType ValueType
	SortOrder int
}

// AttributeValue — значение атрибута для конкретного товара (резолвленное:
// внутри лежит сам Attribute, а не его ID). Используется в доменных слоях.
type AttributeValue struct {
	Attribute Attribute
	ValueText *string
	ValueInt  *int
}

// AttributeValueRow — сырая связка product/attribute/value. Провайдер
// отдаёт это; usecase резолвит AttributeID → Attribute и собирает AttributeValue.
type AttributeValueRow struct {
	ProductID   int64
	AttributeID int64
	ValueText   *string
	ValueInt    *int
}
