package controller

import (
	"github.com/bbrighter/dreams-api/internal/categories"
	dreamcategories "github.com/bbrighter/dreams-api/internal/dreamCategories"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/statistics"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	handler *gin.Engine,
	d *dreams.DreamsService,
	c *categories.CategoriesService,
	dc *dreamcategories.DreamCategoriesService,
	s *statistics.StatisticsService,
) {
	h := handler.Group("/")
	{
		newDreamsRoute(h, d, dc)
		newCategoriesRoute(h, c)
		newStatisticsRoute(h, s)
		h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
