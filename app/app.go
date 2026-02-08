package app

import (
	"time"

	"github.com/bbrighter/dreams-api/config"
	"github.com/bbrighter/dreams-api/docs"
	"github.com/bbrighter/dreams-api/internal/controller"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/bbrighter/dreams-api/internal/usecase/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run(cfg *config.Config, logger *zap.Logger) {
	db := repository.NewDatabase(cfg.DbName, logger)
	dreamsRepo := repository.NewDreamsRepo(db)
	dreamsUseCase := usecase.NewDreamUseCase(dreamsRepo, logger)
	privateDreamsUseCase := usecase.NewPrivateDreamUseCase(dreamsRepo)
	categoriesRepo := repository.NewCategoriesRepo(db)
	categoriesUseCase := usecase.NewCategoriesUseCase(categoriesRepo, logger)
	statisticsRepo := repository.NewStatisticsRepo(db)
	statisticsUseCase := usecase.NewStatisticsUseCase(statisticsRepo)
	authRepo := repository.NewAuthRepo()
	authUseCase := usecase.NewAuthUseCase(authRepo)
	mgmtRepo := repository.NewManagementRepo(db)
	categoriesManagerUseCase := usecase.NewCategoriesManager(mgmtRepo, statisticsRepo, logger)
	repository.Migration(db, logger)

	handler := gin.New()
	handler.Use(
		cors.New(cors.Config{
			AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			MaxAge:          12 * time.Hour,
			AllowAllOrigins: true,
		}),
	)

	controller.NewRouter(
		handler,
		dreamsUseCase,
		privateDreamsUseCase,
		categoriesUseCase,
		statisticsUseCase,
		authUseCase,
		categoriesManagerUseCase,
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
