package v1

import (
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine,
	d usecase.Dreams,
	pd usecase.PrivateDreams,
	p usecase.PersonsAdderRemover,
	c usecase.CategoriesAdderRemover,
	s usecase.Statistics,
	cl usecase.CategoriesLister,
	pl usecase.PersonsLister,
) {
	h := handler.Group("/v1")
	{
		newDreamsRoute(h, d, p, c)
		newPrivateDreamsRoute(h, pd)
		newPersonsRoute(h, pl)
		newCategoriesRoute(h, cl)
		newStatisticsRoute(h, s)
	}
}
