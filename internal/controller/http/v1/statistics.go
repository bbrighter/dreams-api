package v1

import (
	"net/http"
	"strconv"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type statisticsRoute struct {
	c usecase.Statistics
}

func newStatisticsRoute(handler *gin.RouterGroup, c usecase.Statistics) {
	r := &statisticsRoute{c: c}

	h := handler.Group("/statistics")
	{
		h.GET("", r.GetStatistics)
	}
	p := handler.Group("/private/statistics")
	{
		p.GET("", r.GetPrivateStatistics)
	}
}

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} entity.CountsResponse "Counts by category and persons"
// @Router /statistics [get]
// @Param limit query number false "Limit of returned results"
func (r *statisticsRoute) GetStatistics(g *gin.Context) {
	limitStr, exists := g.GetQuery("limit")
	var limit int = 0
	var err error
	if exists {
		limit, err = strconv.Atoi(limitStr)
	}
	if err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	cats, pers := r.c.GetStatistics(false, limit)
	var resp = entity.CountsResponse{
		Categories: cats.ToResponse(),
		Persons:    pers.ToResponse(),
	}
	g.JSON(http.StatusOK, resp)
}

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} entity.CountsResponse "Counts by category and persons"
// @Router /private/statistics [get]
// @Param limit query number false "Limit of returned results"
func (r *statisticsRoute) GetPrivateStatistics(g *gin.Context) {

	limitStr, exists := g.GetQuery("limit")
	var limit int = 0
	var err error
	if exists {
		limit, err = strconv.Atoi(limitStr)
	}
	if err != nil {
		g.Status(http.StatusBadRequest)
		return
	}

	cats, pers := r.c.GetStatistics(true, limit)
	var resp = entity.CountsResponse{
		Categories: cats.ToResponse(),
		Persons:    pers.ToResponse(),
	}
	g.JSON(http.StatusOK, resp)
}
