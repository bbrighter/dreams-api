package v1

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type categoriesRoute struct {
	c usecase.CategoriesLister
	m usecase.ICategoriesManager
}

func newCategoriesRoute(
	handler *gin.RouterGroup,
	c usecase.CategoriesLister,
	m usecase.ICategoriesManager,
) {
	r := &categoriesRoute{c: c, m: m}

	h := handler.Group("/categories")
	{
		h.GET("", r.GetAll)
		h.PATCH("/:id", r.ChangeType)
		h.DELETE("/:id", r.Delete)
		h.POST("/merge", r.Merge)
	}
}

// @Description Get all categories
// @Produce json
// @Success 200 {object} entity.CategoriesResponse
// @Router /categories [get]
// @Param includes query string false "Comma separated list of child objects. Possible entries: dreamsCount"
func (r *categoriesRoute) GetAll(g *gin.Context) {
	includes := entity.ParseIncludes(g.Query("includes"))
	categories := r.c.List(includes)
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
	err = r.m.Delete(id)
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
// @Router /categories/:id [patch]
// @Param newType query string true "New type for this category"
func (r *categoriesRoute) ChangeType(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if handleError(g, err) {
		return
	}
	newTypeStr, err := parseQueryParamString(g, "newType")
	if handleError(g, err) {
		return
	}
	newType, err := entity.NewCategoryType(newTypeStr)
	if handleError(g, err) {
		return
	}

	err = r.m.ChangeType(id, newType)
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
// @Success 200 {object} entity.Categories
// @Router /categories/merge [post]
// @Param mergeCategoriesParams body MergeCategoriesParams true "Which categories should be merged"
func (r *categoriesRoute) Merge(g *gin.Context) {
	var params MergeCategoriesParams
	if err := g.BindJSON(&params); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	cats, err := r.m.Merge(params.SourceCategoryId, params.TargetCategoryId, params.NewName)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, cats.ToResponse())
}
