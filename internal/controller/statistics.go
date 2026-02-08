package controller

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type statisticsRoute struct {
	c usecase.Statistics
	a usecase.Auth
}

func newStatisticsRoute(handler *gin.RouterGroup, c usecase.Statistics, a usecase.Auth) {
	r := &statisticsRoute{c: c, a: a}

	h := handler.Group("/count-categories")
	{
		h.GET("", r.GetCategoriesCount)
		h.GET("/monthly", r.GetMonthlCount)
	}
}

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} entity.CategoriesCountResponse "Counts by category and persons"
// @Router /count-categories [get]
// @Param limit query number false "Limit of returned results"
func (r *statisticsRoute) GetCategoriesCount(g *gin.Context) {
	limit, err := parseQueryParamInt(g, "limit")
	if handleError(g, err) {
		return
	}

	counts, err := r.c.CountCategories(g.Request.Context(), limit)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, counts.ToResponse())
}

// @Description Get count per month
// @Produce json
// @Success 200 {object} entity.Statistics "Monthly statistics"
// @Router /count-categories/monthly [get]
func (r *statisticsRoute) GetMonthlCount(g *gin.Context) {
	cats, dreams, err := r.c.CountByMonth(g.Request.Context(), true)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, dreams.ToResponse(cats))

}
