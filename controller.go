package main

import "github.com/gin-gonic/gin"

type Controller struct {
	Repo Repo
}

func InitController(repo Repo) Controller {
	return Controller{Repo: repo}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Request-Method", "*")
		c.Writer.Header().Set("Access-Control-Content-Type", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accept, origin, Cache-Control, If-None-Match")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func SetupRouter(con Controller) *gin.Engine {
	router := gin.New()
	router.Use(CORSMiddleware())

	router.GET("/dreams", con.GetDreams)
	router.POST("/dreams", con.CreateDream)
	router.GET("/dreams/:id", con.GetDream)
	router.DELETE("/dreams/:id", con.DeleteDream)
	router.PATCH("dreams/:id", con.UpdateDream)
	return router
}
