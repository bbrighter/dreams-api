package v1

import (
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, d usecase.Dreams, p usecase.Persons, c usecase.Categories, s usecase.Statistics) {
	h := handler.Group("/v1")
	{
		newDreamsRoute(h, d, p, c)
		newPrivateDreamsRoute(h, d)
		newPersonsRoute(h, p)
		newCategoriesRoute(h, c)
		newStatisticsRoute(h, s)
	}
}
