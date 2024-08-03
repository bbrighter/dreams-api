package v1

import (
	"net/http"
	"strconv"

	customerrors "github.com/bbrighter/dreams-api/internal/customErrors"
	"github.com/gin-gonic/gin"
)

func parseParamUint(g *gin.Context, paramName string) (uint, error) {
	str := g.Param(paramName)
	id, err := stringToUint(str)
	if err != nil {
		return 0, g.AbortWithError(http.StatusBadRequest, err)
	}
	return id, nil
}

func parseQueryParamString(g *gin.Context, queryParamName string) (string, error) {
	str := g.Query(queryParamName)
	if str == "" {
		return "", g.AbortWithError(http.StatusBadRequest, customerrors.ErrorParameterMissing(queryParamName))

	}
	return str, nil
}

func stringToUint(s string) (uint, error) {
	ui, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(ui), nil
}
