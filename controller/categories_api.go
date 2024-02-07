package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories" validate:"required"`
}

type CategoryResponse struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// @Description Get all categories
// @Produce json
// @Success 200 {object} CategoriesResponse
// @Router /categories [get]
func (con Controller) GetCategories(g *gin.Context) {
	categories := categoriesToCategoriesResponse(con.Repo.GetCategories())
	g.JSON(http.StatusOK, categories)
}

// @Description Add a category to a dream
// @Produce json
// @Success 200 {object} CategoriesResponse
// @Router /dreams/{dreamId}/categories [put]
// @Param name query string true "Name of a category"
func (con Controller) AddCategory(g *gin.Context) {
	dreamId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	categoryName := g.Query("name")
	if categoryName == "" {
		g.AbortWithError(http.StatusBadRequest, ErrorParameterMissing("name"))
		return
	}

	categories, err := con.Repo.AddCategoryToDream(categoryName, uint(dreamId))
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, categoriesToCategoriesResponse(categories))
}

// @Description Remove a category from a dream
// @Produce json
// @Success 200 {object} CategoriesResponse
// @Router /dreams/{dreamId}/categories/{categoryId} [delete]
func (con Controller) RemoveCategory(g *gin.Context) {
	dreamId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	categoryId, err := strconv.Atoi(g.Param("categoryId"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	categories, err := con.Repo.RemoveCategoryFromDream(uint(categoryId), uint(dreamId))
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, categoriesToCategoriesResponse(categories))
}
