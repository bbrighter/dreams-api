package v1

import (
	"net/http"
	"strconv"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/gin-gonic/gin"
)

func parseParamUint(g *gin.Context, paramName string) (uint, error) {
	str := g.Param(paramName)
	id, err := stringToUint(str)
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return 0, entity.ErrorBadParam
	}
	return id, nil
}

func parseQueryParamString(g *gin.Context, queryParamName string) (string, error) {
	str := g.Query(queryParamName)
	if str == "" {
		g.AbortWithError(http.StatusBadRequest, entity.ErrorBadParam)
		return "", entity.ErrorBadParam
	}
	return str, nil
}

func stringToUint(s string) (uint, error) {
	ui, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, entity.ErrorBadParam
	}
	return uint(ui), nil
}
