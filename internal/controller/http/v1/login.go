package v1

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type loginRoutes struct {
	l usecase.Auth
}

func newLoginRoute(handler *gin.RouterGroup, l usecase.Auth) {
	r := &loginRoutes{l}

	handler.POST("/login", r.Login)
	handler.POST("/logout", r.Logout)
}

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Description Login
// @Produce json
// @Success 200 {object} entity.LoginResponse
// @Failure 400
// @Router /login [post]
// @Param loginRequest body LoginRequest true "Password and user name"
func (r loginRoutes) Login(g *gin.Context) {
	var loginParams LoginRequest
	if err := g.BindJSON(&loginParams); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token := r.l.SetAuth(loginParams.Name, loginParams.Password)
	if token == "" {
		g.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	g.JSON(http.StatusOK, entity.LoginResponse{Token: token})
}

// @Description Logout
// @Produce json
// @Success 200
// @Failure 500
// @Router /logout [post]
func (r loginRoutes) Logout(g *gin.Context) {
	r.l.SetAuth("Benni", "")
	g.Status(http.StatusOK)
}
