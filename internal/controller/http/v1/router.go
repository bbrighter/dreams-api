package v1

import (
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine,
	d usecase.Dreams,
	pd usecase.PrivateDreams,
	c usecase.CategoriesAdderRemover,
	s usecase.Statistics,
	cl usecase.CategoriesLister,
) {
	h := handler.Group("/")
	{
		newDreamsRoute(h, d, c)
		newPrivateDreamsRoute(h, pd)
		newCategoriesRoute(h, cl)
		newStatisticsRoute(h, s)
	}
}
