package v1

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type categoriesRoute struct {
	c usecase.CategoriesLister
}

func newCategoriesRoute(handler *gin.RouterGroup, c usecase.CategoriesLister) {
	r := &categoriesRoute{c: c}

	h := handler.Group("/categories")
	{
		h.GET("", r.GetAll)
	}
}

// @Description Get all categories
// @Produce json
// @Success 200 {object} entity.CategoriesResponse
// @Router /dreams-api/v1/categories [get]
func (r *categoriesRoute) GetAll(g *gin.Context) {
	categories := r.c.List()
	g.JSON(http.StatusOK, categories.ToResponse())
}
