package controller

import (
	"net/http"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
	"github.com/gin-gonic/gin"
)

func handleError(g *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	status := apperrors.HttpStatus(err)
	if status == http.StatusInternalServerError {
		g.AbortWithError(status, err)
		return true
	}
	g.AbortWithStatus(status)
	return true
}
