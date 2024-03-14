package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var connectionKey string = "connection"

func AuthenticationMiddleware() gin.HandlerFunc {
	password := "080388"

	return func(c *gin.Context) {
		_, pw, ok := c.Request.BasicAuth()
		if !ok || pw != password {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

func ConnectToRepository(con Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(connectionKey, con)
	}
}

func GetCon(c *gin.Context) Controller {
	con := c.MustGet(connectionKey).(Controller)
	return con
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Request-Method", "*")
		c.Writer.Header().Set("Access-Control-Content-Type", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accept, origin, Cache-Control, If-None-Match, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, DELETE, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
