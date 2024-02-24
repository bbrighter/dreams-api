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

// @Description Get count per category and person
// @Produce json
// @Success 200 {object} CountsResponse "Counts by category and persons"
// @Router /statistics [get]
func GetCountCategories(g *gin.Context) {
	con := GetCon(g)
	categoryCount := con.Repo.CountCategories(false)
	personCount := con.Repo.CountPersons(false)
	var countsResponse CountsResponse
	countsResponse.Categories = countsToCounts(categoryCount)
	countsResponse.Persons = countsToCounts(personCount)
	g.JSON(http.StatusOK, countsResponse)
}
