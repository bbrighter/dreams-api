package controller

import (
	"flag"

	docs "github.com/bbrighter/dreams-api/docs"
	"github.com/bbrighter/dreams-api/store"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Dreams API

type Controller struct {
	Repo store.Repo
}

func InitController(repo store.Repo) Controller {
	return Controller{Repo: repo}
}

func setupRouter(con Controller) *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(ConnectToRepository(con), CORSMiddleware())

	dreamsGroup := router.Group("/dreams")
	dreamsGroup.GET("", GetDreams)
	dreamsGroup.POST("", CreateDream)
	dreamsGroup.GET("/:id", GetDream)
	dreamsGroup.DELETE("/:id", DeleteDream)
	dreamsGroup.PATCH("/:id", UpdateDream)
	dreamsGroup.PUT("/:id/categories", AddCategory)
	dreamsGroup.DELETE("/:id/categories/:categoryId", RemoveCategory)
	dreamsGroup.PUT("/:id/persons", PutPersonToDream)
	dreamsGroup.DELETE("/:id/persons/:personId", RemovePersonFromDream)

	privateDreams := router.Group("/private/dreams", AuthenticationMiddleware())
	privateDreams.GET("", GetPrivateDreams)
	privateDreams.GET("/:id", GetPrivateDream)
	privateDreams.PATCH("/:id", TogglePrivateDream)

	categoriesGroup := router.Group("/categories")
	categoriesGroup.GET("", GetCategories)

	personsGroup := router.Group("/persons")
	personsGroup.GET("", GetPersons)

	statisticsGroup := router.Group("/statistics")
	statisticsGroup.GET("", GetCountCategories)
	privateStatistics := router.Group("/private/statistics", AuthenticationMiddleware())
	privateStatistics.GET("", GetPrivateCountCategories)
	return router
}

func (con Controller) RunRouter() {
	var ipAddress string
	flag.StringVar(&ipAddress, "ipAddress", "localhost", "IP address to run")
	flag.Parse()

	var host = ipAddress + ":5005"
	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "2.0"

	var router *gin.Engine = setupRouter(con)

	router.Run(host)
}
