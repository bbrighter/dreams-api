package v1

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type personsRoute struct {
	p usecase.Persons
}

func newPersonsRoute(handler *gin.RouterGroup, p usecase.Persons) {
	r := &personsRoute{p}

	h := handler.Group("/persons")
	{
		h.GET("", r.GetAll)
	}
}

// @Description Get all persons
// @Produce json
// @Success 200 {object} entity.PersonsResponse
// @Router /v1/persons [get]
func (r *personsRoute) GetAll(g *gin.Context) {
	persons := r.p.GetAll()
	g.JSON(http.StatusOK, persons.ToResponse())
}
