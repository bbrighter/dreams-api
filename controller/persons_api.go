package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PersonResponse struct {
	ID   uint   `json:"id"  validate:"required"`
	Name string `json:"name"  validate:"required"`
}

type PersonsResponse struct {
	Persons []PersonResponse `json:"persons" validate:"required"`
}

// @Description Get all persons
// @Produce json
// @Success 200 {object} []PersonResponse
// @Router /persons [get]
func (con Controller) GetPersons(g *gin.Context) {
	persons := con.Repo.GetAllPersons()
	var personsResp []PersonResponse = []PersonResponse{}
	for _, p := range persons {
		personsResp = append(personsResp, personToPersonResponse(p))
	}
	g.JSON(http.StatusOK, personsResp)
}

// @Description Add a person to a dream
// @Produce json
// @Success 200 {number} id
// @Router /dreams/{dreamId}/persons [put]
// @Param name query string true "Name of person"
func (con Controller) PutPersonToDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	name, err := parseQueryParamString(g, "name")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := con.Repo.AddPersonToDream(name, dreamId)
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, id)
}

// @Description Delete a person from a dream
// @Produce json
// @Success 200 {object} PersonsResponse
// @Router /dreams/{dreamId}/persons/{personId} [delete]
func (con Controller) RemovePersonFromDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	personId, err := parseParamUint(g, "personId")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	persons, err := con.Repo.RemovePersonFromDream(personId, dreamId)
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, personsToPersonsResponse(persons))
}
