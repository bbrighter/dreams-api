package usecase

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthUseCase struct {
	r IAuthRepo
}

func NewAuthUseCase(r IAuthRepo) AuthUseCase {
	return AuthUseCase{r: r}
}

func (uc AuthUseCase) IsValidToken(token string) bool {
	return uc.r.IsValidToken(token)
}

func (uc AuthUseCase) SetAuth(name, pw string) string {
	return uc.r.SetAuth(name, pw)
}

func (uc AuthUseCase) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !uc.r.IsValidToken(auth) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
