package controller

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/statistics"
	"github.com/gin-gonic/gin"
)

type statisticsRoute struct {
	c *statistics.StatisticsService
}

func newStatisticsRoute(
	handler *gin.RouterGroup,
	c *statistics.StatisticsService,
) {
	r := &statisticsRoute{c: c}

	h := handler.Group("/count-categories")
	{
		h.GET("/monthly", r.GetMonthlyCount)
	}
}

type Statistics struct {
	Statistics []Statistic `json:"statistics"`
}

type Statistic struct {
	Month      string       `json:"month"`
	Categories []CountByCat `json:"categories"`
	DreamCount int64        `json:"dreamCount"`
}

type CountByCat struct {
	CategoryId uint  `json:"categoryId"`
	Count      int64 `json:"count"`
}

func countsToStatistics(dreams []statistics.CountByMonth, cats []statistics.CountByCatAndMonth) Statistics {
	var statistics = []Statistic{}
	for _, d := range dreams {
		var dreamCats = []CountByCat{}
		for _, c := range cats {
			dreamCats = append(dreamCats, CountByCat{CategoryId: c.CategoryId, Count: c.Count})
		}
		statistics = append(statistics, Statistic{
			Month:      d.Month,
			DreamCount: d.Count,
			Categories: dreamCats,
		})
	}
	return Statistics{Statistics: statistics}
}

// @Description Get count per month
// @Produce json
// @Success 200 {object} Statistics "Monthly statistics"
// @Router /count-categories/monthly [get]
func (r *statisticsRoute) GetMonthlyCount(g *gin.Context) {
	cats, dreams, err := r.c.CountByMonth(g.Request.Context())
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, countsToStatistics(dreams, cats))

}
