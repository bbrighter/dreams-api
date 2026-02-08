package controller

import (
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(handler *gin.Engine,
	d usecase.IDreams,
	pd usecase.PrivateDreams,
	c usecase.ICategories,
	s usecase.Statistics,
	l usecase.Auth,
	cm usecase.ICategoriesManager,
) {
	h := handler.Group("/")
	{
		newDreamsRoute(h, d, c)
		newPrivateDreamsRoute(h, pd, l)
		newCategoriesRoute(h, cm, c)
		newStatisticsRoute(h, s, l)
		newLoginRoute(h, l)
		h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
