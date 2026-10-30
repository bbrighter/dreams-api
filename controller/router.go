package controller

import (
	"github.com/bbrighter/dreams-api/internal/categories"
	dreamcategories "github.com/bbrighter/dreams-api/internal/dreamCategories"
	dt "github.com/bbrighter/dreams-api/internal/dreamTranscriptions"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/recordings"
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
	dts *dt.DreamTranscriptionService,
	rec *recordings.RecordingsService,
	mimeType string,
) {
	h := handler.Group("/")
	{
		newDreamsRoute(h, d, dc)
		newCategoriesRoute(h, c)
		newStatisticsRoute(h, s)
		newRecordingRoute(h, dts, rec, mimeType)
		h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
