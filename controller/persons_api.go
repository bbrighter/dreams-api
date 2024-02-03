package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PersonResponse struct {
	ID   uint
	Name string
}

// @Description Add a person to a dream
// Produce json
// @Success 200
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

	err = con.Repo.AddPersonToDream(name, dreamId)
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.Status(http.StatusOK)
}

// @Description Delete a person from a dream
// Produce json
// @Success 200
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

	err = con.Repo.RemovePersonFromDream(personId, dreamId)
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.Status(200)
}
