package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Count struct {
	ID    uint `json:"id" validate:"required"`
	Count int  `json:"count" validate:"required"`
}

type CountsResponse struct {
	Categories []Count `json:"categories" validate:"required"`
	Persons    []Count `json:"persons" validate:"required"`
}

func (con Controller) getCountCategories(showAll bool, limit int) CountsResponse {
	categoryCount := con.Repo.CountCategories(showAll, limit)
	personCount := con.Repo.CountPersons(showAll, limit)
	var countsResponse CountsResponse
	countsResponse.Categories = countsToCounts(categoryCount)
	countsResponse.Persons = countsToCounts(personCount)
	return countsResponse
}

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} CountsResponse "Counts by category and persons"
// @Router /statistics [get]
// @Param limit query number false "Limit of returned results"
func GetCountCategories(g *gin.Context) {
	con := GetCon(g)
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

	var countsResponse CountsResponse = con.getCountCategories(false, limit)
	g.JSON(http.StatusOK, countsResponse)
}

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} CountsResponse "Counts by category and persons"
// @Failure 401
// @Security BasicAuth
// @Router /private/statistics [get]
// @Param limit query number false "Limit of returned results"
func GetPrivateCountCategories(g *gin.Context) {
	con := GetCon(g)

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
	var countsResponse CountsResponse = con.getCountCategories(true, limit)
	g.JSON(http.StatusOK, countsResponse)
}
