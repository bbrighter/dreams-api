package app

import (
	"time"

	"github.com/bbrighter/dreams-api/config"
	"github.com/bbrighter/dreams-api/docs"
	v1 "github.com/bbrighter/dreams-api/internal/controller/http/v1"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/bbrighter/dreams-api/internal/usecase/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run(cfg *config.Config, logger *zap.Logger) {
	db := repository.NewDatabase(cfg.DbName, logger)
	dreamsRepo := repository.NewDreamsRepo(db)
	dreamsUseCase := usecase.NewDreamUseCase(dreamsRepo)
	privateDreamsUseCase := usecase.NewPrivateDreamUseCase(dreamsRepo)
	personsRepo := repository.NewPersonsRepo(db)
	personsUseCase := usecase.NewPersonsUseCase(personsRepo)
	categoriesRepo := repository.NewCategoriesRepo(db)
	categoriesUseCase := usecase.NewCategoriesUseCase(categoriesRepo)
	statisticsRepo := repository.NewStatisticsRepo(db)
	statisticsUseCase := usecase.NewStatisticsUseCase(statisticsRepo)
	migration(db, logger)

	handler := gin.New()
	handler.Use(
		cors.New(cors.Config{
			AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			MaxAge:          12 * time.Hour,
			AllowAllOrigins: true,
		}),
	)

	v1.NewRouter(
		handler,
		dreamsUseCase,
		privateDreamsUseCase,
		personsUseCase,
		categoriesUseCase,
		statisticsUseCase,
		categoriesUseCase,
		personsUseCase,
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
