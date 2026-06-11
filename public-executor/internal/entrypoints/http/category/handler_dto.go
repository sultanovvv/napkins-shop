package category

import (
	"napkins-shop/public-executor/api"
	"shared/entity"
)

func toAPICategoryNode(c entity.Category) api.CategoryNode {
	node := api.CategoryNode{
		Id:        c.ID,
		Slug:      c.Slug,
		Name:      c.Name,
		SortOrder: c.SortOrder,
		Children:  make([]api.CategoryNode, len(c.Children)),
	}
	for i, child := range c.Children {
		node.Children[i] = toAPICategoryNode(child)
	}
	return node
}
