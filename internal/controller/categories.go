package controller

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type categoriesRoute struct {
	m usecase.ICategoriesManager
	c usecase.ICategories
}

func newCategoriesRoute(
	handler *gin.RouterGroup,
	m usecase.ICategoriesManager,
	c usecase.ICategories,
) {
	r := &categoriesRoute{m: m, c: c}

	h := handler.Group("/categories")
	{
		h.GET("", r.List)
		h.POST("/merge", r.Merge)
		idGroup := h.Group("/:id")
		{
			idGroup.PATCH("/type", r.ChangeType)
			idGroup.PATCH("/name", r.ChangeName)
			idGroup.DELETE("", r.Delete)
		}
	}
}

// @Description Get all categories
// @Produce json
// @Success 200 {object} entity.CategoriesResponse
// @Router /categories [get]
func (r *categoriesRoute) List(g *gin.Context) {
	categories, err := r.c.List(g.Request.Context())
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, categories.ToResponse())
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
	err = r.m.Delete(g.Request.Context(), id)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
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
	newType, err := entity.NewCategoryType(newTypeStr)
	if handleError(g, err) {
		return
	}

	err = r.m.ChangeType(g.Request.Context(), id, newType)
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

	err = r.m.ChangeName(g.Request.Context(), id, newName)
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
// @Success 200 {object} entity.CategoriesCountResponse
// @Router /categories/merge [post]
// @Param mergeCategoriesParams body MergeCategoriesParams true "Which categories should be merged"
func (r *categoriesRoute) Merge(g *gin.Context) {
	var params MergeCategoriesParams
	if err := g.BindJSON(&params); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	cats, err := r.m.Merge(g.Request.Context(), params.SourceCategoryId, params.TargetCategoryId, params.NewName)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, cats.ToResponse())
}
