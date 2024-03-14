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
func GetCategories(g *gin.Context) {
	con := GetCon(g)
	categories := categoriesToCategoriesResponse(con.Repo.GetCategories())
	g.JSON(http.StatusOK, categories)
}

// @Description Add a category to a dream
// @Produce json
// @Success 200 {object} CategoriesResponse
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/categories [put]
// @Param name query string true "Name of a category"
func AddCategory(g *gin.Context) {
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

	con := GetCon(g)
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
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/categories/{categoryId} [delete]
func RemoveCategory(g *gin.Context) {
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

	con := GetCon(g)
	categories, err := con.Repo.RemoveCategoryFromDream(uint(categoryId), uint(dreamId))
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, categoriesToCategoriesResponse(categories))
}
