package app

import (
	"time"

	"github.com/bbrighter/dreams-api/config"
	"github.com/bbrighter/dreams-api/docs"
	v1 "github.com/bbrighter/dreams-api/internal/controller/http/v1"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/bbrighter/dreams-api/internal/usecase/repository"
	logger "github.com/bbrighter/zapLogWrapper"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	opts := logger.NewLoggerOptions()
	if cfg.Logs.Folder != "" {
		opts.SetFolder(cfg.Logs.Folder)
	}
	log := logger.NewLogger(opts)

	repo := repository.NewDreamsRepo(cfg.DB.Name, log)
	dreamsUseCase := usecase.New(repo)
	personsRepo := repository.NewPersonsRepo(cfg.DB.Name, log)
	personsUseCase := usecase.NewPersonsUseCase(personsRepo)
	categoriesRepo := repository.NewCategoriesRepo(cfg.DB.Name, log)
	categoriesUseCase := usecase.NewCategoriesUseCase(categoriesRepo)
	statisticsRepo := repository.NewStatisticsRepo(cfg.DB.Name, log)
	statisticsUseCase := usecase.NewStatisticsUseCase(statisticsRepo)
	Migration(repo.Repo)

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
		personsUseCase,
		categoriesUseCase,
		statisticsUseCase,
	)

	var host = cfg.API.Host + ":" + cfg.API.Port
	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "2.0"

	log.Info("Starting API", zap.String("Port", cfg.API.Port), zap.String("Host", cfg.API.Host))
	err := handler.Run(host)
	if err != nil {
		log.Fatal("Cannot start API", zap.Error(err))
	}
}
