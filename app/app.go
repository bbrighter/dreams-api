package app

import (
	"github.com/bbrighter/dreams-api/config"
	"github.com/bbrighter/dreams-api/docs"
	v1 "github.com/bbrighter/dreams-api/internal/controller/http/v1"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/bbrighter/dreams-api/internal/usecase/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	repo := repository.NewDreamsRepo(cfg.DB.Name)
	dreamsUseCase := usecase.New(repo)
	personsRepo := repository.NewPersonsRepo(cfg.DB.Name)
	personsUseCase := usecase.NewPersonsUseCase(personsRepo)
	categoriesRepo := repository.NewCategoriesRepo(cfg.DB.Name)
	categoriesUseCase := usecase.NewCategoriesUseCase(categoriesRepo)
	statisticsRepo := repository.NewStatisticsRepo(cfg.DB.Name)
	statisticsUseCase := usecase.NewStatisticsUseCase(statisticsRepo)
	Migration(repo.Repo)

	handler := gin.New()
	handler.Use(cors.New(cors.Config{
		// AllowAllOrigins: true,
		AllowMethods:  []string{"GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowOrigins:  []string{"http://" + cfg.API.Host + ":3005", "http://localhost:*"},
		AllowWildcard: true,
	}))

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

	handler.Run(host)
}
