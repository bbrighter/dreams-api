package controller

import (
	docs "github.com/bbrighter/dreams-api/docs"
	"github.com/bbrighter/dreams-api/store"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Dreams API
// @version 2.0
// @BasePath /

type Controller struct {
	Repo store.Repo
}

func InitController(repo store.Repo) Controller {
	return Controller{Repo: repo}
}

func SetupRouter(con Controller) *gin.Engine {
	router := gin.New()
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "1.0"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(CORSMiddleware(), ConnectToRepository(con))

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

	privateGroup := router.Group("/private/dreams")
	privateGroup.Use(AuthenticationMiddleware())
	privateGroup.GET("", GetPrivateDreams)
	privateGroup.GET("/:id", GetPrivateDream)
	privateGroup.PATCH("/:id", TogglePrivateDream)

	categoriesGroup := router.Group("/categories")
	categoriesGroup.GET("", GetCategories)

	personsGroup := router.Group("/persons")
	personsGroup.GET("", GetPersons)

	statisticsGroup := router.Group("/statistics")
	statisticsGroup.GET("", GetCountCategories)
	return router
}
