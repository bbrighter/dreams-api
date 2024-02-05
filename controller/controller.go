package controller

import (
	docs "github.com/bbrighter/dreams-api/docs"
	"github.com/bbrighter/dreams-api/store"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Dreams API
// @version 1.0
// @BasePath /

type Controller struct {
	Repo store.Repo
}

func InitController(repo store.Repo) Controller {
	return Controller{Repo: repo}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Request-Method", "*")
		c.Writer.Header().Set("Access-Control-Content-Type", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accept, origin, Cache-Control, If-None-Match")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, DELETE, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func SetupRouter(con Controller) *gin.Engine {
	router := gin.New()
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "1.0"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(CORSMiddleware())

	dreamsGroup := router.Group("/dreams")
	dreamsGroup.GET("", con.GetDreams)
	dreamsGroup.POST("", con.CreateDream)
	dreamsGroup.GET("/:id", con.GetDream)
	dreamsGroup.DELETE("/:id", con.DeleteDream)
	dreamsGroup.PATCH("/:id", con.UpdateDream)
	dreamsGroup.PUT("/:id/tags", con.AddTag)
	dreamsGroup.DELETE("/:id/tags/:tagId", con.RemoveTag)
	dreamsGroup.PUT("/:id/persons", con.PutPersonToDream)
	dreamsGroup.DELETE("/:id/persons/:personId", con.RemovePersonFromDream)

	tagsGroup := router.Group("/tags")
	tagsGroup.GET("", con.GetTags)

	personsGroup := router.Group("/persons")
	personsGroup.GET("", con.GetPersons)
	return router
}
