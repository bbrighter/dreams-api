package controller

import (
	"errors"
	"net/http"

	"github.com/bbrighter/dreams-api/internal/categories"
	"github.com/bbrighter/dreams-api/internal/entities"
	"github.com/gin-gonic/gin"
)

type categoriesRoute struct {
	c *categories.CategoriesService
}

func newCategoriesRoute(
	handler *gin.RouterGroup,
	c *categories.CategoriesService,
) {
	r := &categoriesRoute{c: c}

	h := handler.Group("/categories")
	{
		h.GET("", r.List)
		h.GET("/with-count", r.ListWithCount)
		h.POST("/merge", r.Merge)
		idGroup := h.Group("/:id")
		{
			idGroup.PATCH("/type", r.ChangeType)
			idGroup.PATCH("/name", r.ChangeName)
			idGroup.DELETE("", r.Delete)
		}
	}
}

type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

type CategoryResponse struct {
	ID    uint   `json:"id" `
	Name  string `json:"name"`
	Type  string `json:"type"`
	Count int64  `json:"count,omitempty"`
}

func categoriesToResponse(cats []entities.Category) CategoryListResponse {
	var resp = []CategoryResponse{}
	for _, cat := range cats {
		resp = append(resp, CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
			Type: string(cat.Type),
		})
	}
	return CategoryListResponse{Categories: resp}
}

func countCategoriesToResponse(cats []categories.CountByCategory) CategoryListResponse {
	var resp = []CategoryResponse{}
	for _, cat := range cats {
		resp = append(resp, CategoryResponse{
			ID:    cat.ID,
			Name:  cat.Name,
			Type:  string(cat.Type),
			Count: cat.Count,
		})
	}
	return CategoryListResponse{Categories: resp}
}

// @Description Get all categories
// @Produce json
// @Success 200 {object} CategoryListResponse
// @Router /categories [get]
func (r *categoriesRoute) List(g *gin.Context) {
	categories, err := r.c.ListCategories(g.Request.Context())
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, categoriesToResponse(categories))
}

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} CategoryListResponse "Counts by category and persons"
// @Router /categories/with-count [get]
func (r *categoriesRoute) ListWithCount(g *gin.Context) {
	counts, err := r.c.ListAndCountCategories(g.Request.Context())
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, countCategoriesToResponse(counts))
}

// @Description Delete a category. Must be contained in no dreams.
// @Produce json
// @Success 200
// @Error 400
// @Error 404
// @Router /categories/:id [delete]
func (r *categoriesRoute) Delete(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if handleError(g, err) {
		return
	}
	err = r.c.DeleteCategory(g.Request.Context(), id)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
}

var validCategoryTypes = map[string]bool{
	"person":   true,
	"category": true,
}

// @Description Change the type of a category.
// @Produce json
// @Success 200
// @Error 400
// @Error 404
// @Router /categories/:id/type [patch]
// @Param type query string true "New type for this category"
func (r *categoriesRoute) ChangeType(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if handleError(g, err) {
		return
	}
	newTypeStr, err := parseQueryParamString(g, "type")
	if handleError(g, err) {
		return
	}

	if !validCategoryTypes[newTypeStr] {
		handleError(g, errors.New("Invalid type"))
		return
	}
	newType := entities.CategoryType(newTypeStr)

	err = r.c.UpdateCategory(g.Request.Context(), id, nil, &newType)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
}

// @Description Change the name of a category.
// @Produce json
// @Success 200
// @Error 400
// @Error 404
// @Router /categories/:id/name [patch]
// @Param name query string true "New name for this category"
func (r *categoriesRoute) ChangeName(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if handleError(g, err) {
		return
	}
	newName, err := parseQueryParamString(g, "name")
	if handleError(g, err) {
		return
	}

	err = r.c.UpdateCategory(g.Request.Context(), id, &newName, nil)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
}

type MergeCategoriesParams struct {
	SourceCategoryId uint   `json:"sourceCategoryId" binding:"required"`
	TargetCategoryId uint   `json:"targetCategoryId" binding:"required"`
	NewName          string `json:"newName" binding:"required"`
}

// @Description Merge two categories.
// @Produce json
// @Success 200 {object} CategoryListResponse
// @Router /categories/merge [post]
// @Param mergeCategoriesParams body MergeCategoriesParams true "Which categories should be merged"
func (r *categoriesRoute) Merge(g *gin.Context) {
	var params MergeCategoriesParams
	if err := g.BindJSON(&params); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err := r.c.MergeCategories(g.Request.Context(), params.SourceCategoryId, params.TargetCategoryId, params.NewName)
	if handleError(g, err) {
		return
	}
	cats, err := r.c.ListAndCountCategories(g.Request.Context())
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, countCategoriesToResponse(cats))
}
