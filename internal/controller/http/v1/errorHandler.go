package v1

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/gin-gonic/gin"
)

func handleError(g *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if customError, ok := err.(*entity.CustomError); ok {
		if customError.IsNotFound() {
			g.AbortWithStatus(http.StatusNotFound)
			return true
		}
		if customError.IsBadParam() {
			g.AbortWithError(http.StatusBadRequest, err)
			return true
		}
	}
	g.AbortWithError(http.StatusInternalServerError, err)
	return true
}
