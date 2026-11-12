package category

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"napkins-shop/backoffice-executor/api"
	categoryuc "napkins-shop/backoffice-executor/internal/usecases/category"
)

type CategoryHandler struct {
	usecase categoryuc.ICategoryUseCases
}

func NewCategoryHandler(usecase categoryuc.ICategoryUseCases) *CategoryHandler {
	return &CategoryHandler{usecase: usecase}
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var body api.CreateCategoryRequest
	fields, err := readJSON(c, &body, "name", "slug", "parentId", "sortOrder")
	if err != nil {
		return err
	}
	if isNull(fields["sortOrder"]) {
		return invalidField("sortOrder", "must not be null", nil)
	}
	in := categoryuc.CreateCategoryInUDTO{Name: body.Name, Slug: body.Slug, ParentID: body.ParentId}
	if body.SortOrder != nil {
		in.SortOrder = *body.SortOrder
	}
	created, err := h.usecase.Create(c.Request().Context(), in)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toAPI(*created))
}

func (h *CategoryHandler) ListCategories(c echo.Context) error {
	items, err := h.usecase.List(c.Request().Context())
	if err != nil {
		return err
	}
	result := make([]api.Category, len(items))
	for i, item := range items {
		result[i] = api.Category{
			Id: item.ID, Name: item.Name, Slug: item.Slug,
			ParentId: item.ParentID, SortOrder: item.SortOrder,
		}
	}
	return c.JSON(http.StatusOK, result)
}

func (h *CategoryHandler) GetCategory(c echo.Context, id string) error {
	categoryID, err := parseID(id)
	if err != nil {
		return err
	}
	item, err := h.usecase.GetByID(c.Request().Context(), categoryID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toAPI(*item))
}

func (h *CategoryHandler) UpdateCategory(c echo.Context, id string) error {
	categoryID, err := parseID(id)
	if err != nil {
		return err
	}
	var body api.UpdateCategoryRequest
	fields, err := readJSON(c, &body, "name", "slug", "parentId", "sortOrder", "clearParentId")
	if err != nil {
		return err
	}
	if isNull(fields["clearParentId"]) {
		return invalidField("clearParentId", "must not be null", nil)
	}
	_, nameSet := fields["name"]
	_, slugSet := fields["slug"]
	_, sortOrderSet := fields["sortOrder"]
	in := categoryuc.UpdateCategoryInUDTO{
		ID: categoryID, ParentID: body.ParentId,
		Name:      categoryuc.PatchField[string]{Set: nameSet, Value: body.Name},
		Slug:      categoryuc.PatchField[string]{Set: slugSet, Value: body.Slug},
		SortOrder: categoryuc.PatchField[int]{Set: sortOrderSet, Value: body.SortOrder},
	}
	if body.ClearParentId != nil {
		in.ClearParentID = *body.ClearParentId
	}
	updated, err := h.usecase.Update(c.Request().Context(), in)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toAPI(*updated))
}

func (h *CategoryHandler) DeleteCategory(c echo.Context, id string) error {
	categoryID, err := parseID(id)
	if err != nil {
		return err
	}
	if err := h.usecase.Delete(c.Request().Context(), categoryID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
