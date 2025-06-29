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
	l usecase.Auth,
	cm usecase.ICategoriesManager,
) {
	h := handler.Group("/")
	{
		newDreamsRoute(h, d, c)
		newPrivateDreamsRoute(h, pd, l)
		newCategoriesRoute(h, cl, cm)
		newStatisticsRoute(h, s, l)
		newLoginRoute(h, l)
	}
}
