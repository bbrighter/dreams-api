package v1

import (
	"errors"
	"net/http"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/gin-gonic/gin"
)

func handleError(g *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, entity.ErrorNotFound) {
		g.AbortWithStatus(http.StatusNotFound)
		return true
	}
	if errors.Is(err, entity.ErrorBadParam) || errors.Is(err, entity.ErrorMissingParam) {
		g.AbortWithError(http.StatusBadRequest, err)
		return true
	}
	g.AbortWithError(http.StatusInternalServerError, err)
	return true

}
