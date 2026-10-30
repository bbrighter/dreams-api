package app

import (
	"time"

	"github.com/bbrighter/dreams-api/config"
	"github.com/bbrighter/dreams-api/controller"
	"github.com/bbrighter/dreams-api/docs"
	"github.com/bbrighter/dreams-api/internal/blob"
	"github.com/bbrighter/dreams-api/internal/categories"
	dreamcategories "github.com/bbrighter/dreams-api/internal/dreamCategories"
	dreamtranscriptions "github.com/bbrighter/dreams-api/internal/dreamTranscriptions"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/migrations"
	"github.com/bbrighter/dreams-api/internal/recordings"
	"github.com/bbrighter/dreams-api/internal/statistics"
	"github.com/bbrighter/dreams-api/internal/transcription"
	unitofwork "github.com/bbrighter/dreams-api/internal/unitOfWork"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run(cfg *config.Config, logger *zap.Logger) {
	mimeType := "audio/webm;codecs=opus"
	db := migrations.NewDatabase(cfg.DbName, cfg.DataFolder, logger)
	migrations.Migration(db, logger)
	dreamsRepo := dreams.NewDreamsRepo(db)
	categoriesRepo := categories.NewCategoriesRepo(db)
	statisticsRepo := statistics.NewStatisticsRepo(db)
	transRepo := transcription.NewTranscriptionRepo(cfg.Transcription.Url, cfg.Transcription.Model, mimeType)
	recordRepo := recordings.NewRecordingsRepository(db)
	uow := unitofwork.NewUnitOfWork(db)
	blobManager := blob.NewBlobStorage(cfg.DataFolder)

	handler := gin.New()
	handler.Use(
		cors.New(cors.Config{
			AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			MaxAge:          12 * time.Hour,
			AllowAllOrigins: true,
		}),
	)

	transService := transcription.NewTranscriptionService(transRepo)
	dreamService := dreams.NewDreamsService(dreamsRepo)
	recordingService := recordings.NewRecordingsService(recordRepo, blobManager, mimeType)

	if err := transService.Init(); err != nil {
		panic(err)
	}

	controller.NewRouter(
		handler,
		dreamService,
		categories.NewCategoriesService(categoriesRepo),
		dreamcategories.NewDreamCategoriesService(dreamsRepo, categoriesRepo, uow),
		statistics.NewStatisticsService(statisticsRepo),
		dreamtranscriptions.NewDreamTranscriptionService(transService, dreamService, recordingService),
		recordingService,
		mimeType,
	)

	var host = cfg.Host + ":" + cfg.Port
	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "2.0"

	logger.Info("Starting API", zap.String("Port", cfg.Port), zap.String("Host", cfg.Host))
	err := handler.Run(host)
	if err != nil {
		logger.Fatal("Cannot start API", zap.Error(err))
	}
}
