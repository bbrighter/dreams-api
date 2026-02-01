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
	}
	// p := handler.Group("/private/statistics")
	// {
	// 	p.GET("", r.GetPrivateStatistics, a.AuthMiddleware())
	// }
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

// // @Description Get count per category and person
// // @Produce json
// // @Success 200 {object} entity.CountsResponse "Counts by category and persons"
// // @Router /private/statistics [get]
// // @Security BasicAuth
// // @Param limit query number false "Limit of returned results"
// func (r *statisticsRoute) GetPrivateStatistics(g *gin.Context) {

// 	limitStr, exists := g.GetQuery("limit")
// 	var limit int = 0
// 	var err error
// 	if exists {
// 		limit, err = strconv.Atoi(limitStr)
// 	}
// 	if err != nil {
// 		g.Status(http.StatusBadRequest)
// 		return
// 	}

// 	cats, pers := r.c.GetStatistics(true, limit)
// 	var resp = entity.CategoriesCountResponse{
// 		Categories: cats.ToResponse(),
// 		Persons:    pers.ToResponse(),
// 	}
// 	g.JSON(http.StatusOK, resp)
// }
