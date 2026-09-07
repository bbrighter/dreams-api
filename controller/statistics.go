package controller

import (
	"cmp"
	"net/http"
	"slices"

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

func countsToStatistics(counts []statistics.CountByCatAndMonth) Statistics {
	var monthlyStats = make(map[string]Statistic)
	for _, c := range counts {
		val, ok := monthlyStats[c.Month]
		if !ok {
			val = Statistic{Categories: []CountByCat{}, Month: c.Month}
		}

		if c.ResultType == "total" {
			val.DreamCount = c.Count
		} else if c.CategoryId != nil {
			val.Categories = append(val.Categories,
				CountByCat{
					CategoryId: *c.CategoryId,
					Count:      c.Count,
				})
		}

		monthlyStats[c.Month] = val

	}

	var statistics = make([]Statistic, 0, len(monthlyStats))
	for _, val := range monthlyStats {
		statistics = append(statistics, val)
	}
	slices.SortFunc(statistics, func(a, b Statistic) int {
		return cmp.Compare(a.Month, b.Month)
	})

	return Statistics{Statistics: statistics}
}

// @Description Get count per month
// @Produce json
// @Success 200 {object} Statistics "Monthly statistics"
// @Router /count-categories/monthly [get]
func (r *statisticsRoute) GetMonthlyCount(g *gin.Context) {
	cats, err := r.c.CountByMonth(g.Request.Context())
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, countsToStatistics(cats))

}
