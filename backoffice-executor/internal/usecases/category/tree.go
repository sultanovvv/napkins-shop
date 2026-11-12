package category

import (
	"fmt"
	"sort"

	"napkins-shop/backoffice-executor/internal/customerr"
	"shared/entity"
)

// orderedCategories строит обход «родитель перед детьми» без рекурсии.
// Порядок входных строк и глубина дерева не влияют на сохранность потомков.
func orderedCategories(rows []entity.Category) ([]CategoryListItemOutUDTO, error) {
	byID := make(map[int64]int, len(rows))
	children := make(map[int64][]int, len(rows))
	roots := make([]int, 0)
	for i, row := range rows {
		if row.ID <= 0 {
			return nil, customerr.NewInternal("Internal Server Error", fmt.Errorf("invalid category hierarchy: invalid id %d", row.ID))
		}
		if _, exists := byID[row.ID]; exists {
			return nil, customerr.NewInternal("Internal Server Error", fmt.Errorf("invalid category hierarchy: duplicate id %d", row.ID))
		}
		byID[row.ID] = i
	}
	for i, row := range rows {
		if row.ParentID == nil {
			roots = append(roots, i)
			continue
		}
		if _, exists := byID[*row.ParentID]; !exists {
			return nil, customerr.NewInternal("Internal Server Error", fmt.Errorf("invalid category hierarchy: parent %d of category %d is missing", *row.ParentID, row.ID))
		}
		children[*row.ParentID] = append(children[*row.ParentID], i)
	}
	order := func(indices []int) {
		sort.Slice(indices, func(i, j int) bool {
			left, right := rows[indices[i]], rows[indices[j]]
			if left.SortOrder != right.SortOrder {
				return left.SortOrder < right.SortOrder
			}
			if left.Name != right.Name {
				return left.Name < right.Name
			}
			return left.ID < right.ID
		})
	}
	order(roots)
	for _, indices := range children {
		order(indices)
	}

	stack := make([]int, 0, len(rows))
	for i := len(roots) - 1; i >= 0; i-- {
		stack = append(stack, roots[i])
	}
	out := make([]CategoryListItemOutUDTO, 0, len(rows))
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		row := rows[current]
		out = append(out, CategoryListItemOutUDTO{
			ID: row.ID, Name: row.Name, Slug: row.Slug,
			ParentID: row.ParentID, SortOrder: row.SortOrder,
		})
		descendants := children[row.ID]
		for i := len(descendants) - 1; i >= 0; i-- {
			stack = append(stack, descendants[i])
		}
	}
	// При единственном parent_id и уникальных ID цикл не достижим из
	// корня: узлы циклической компоненты останутся вне результата.
	if len(out) != len(rows) {
		return nil, customerr.NewInternal("Internal Server Error", fmt.Errorf("invalid category hierarchy: cycle detected"))
	}
	return out, nil
}
