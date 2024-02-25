package controller

import (
	"net/http"

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

func (con Controller) getCountCategories(showAll bool) CountsResponse {
	categoryCount := con.Repo.CountCategories(showAll)
	personCount := con.Repo.CountPersons(showAll)
	var countsResponse CountsResponse
	countsResponse.Categories = countsToCounts(categoryCount)
	countsResponse.Persons = countsToCounts(personCount)
	return countsResponse
}

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} CountsResponse "Counts by category and persons"
// @Router /statistics [get]
func GetCountCategories(g *gin.Context) {
	con := GetCon(g)
	var countsResponse CountsResponse = con.getCountCategories(false)
	g.JSON(http.StatusOK, countsResponse)
}

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} CountsResponse "Counts by category and persons"
// @Failure 401
// @Security BasicAuth
// @Router /private/statistics [get]
func GetPrivateCountCategories(g *gin.Context) {
	con := GetCon(g)
	var countsResponse CountsResponse = con.getCountCategories(true)
	g.JSON(http.StatusOK, countsResponse)
}
