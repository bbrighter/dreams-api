package controller

import (
	"strconv"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
	"github.com/gin-gonic/gin"
)

func parseParamUint(g *gin.Context, paramName string) (uint, error) {
	str := g.Param(paramName)
	id, err := stringToUint(str)
	if err != nil {
		return 0, apperrors.ErrBadParam
	}
	return id, nil
}

func parseQueryParamString(g *gin.Context, queryParamName string) (string, error) {
	str := g.Query(queryParamName)
	if str == "" {
		return "", apperrors.ErrBadParam
	}
	return str, nil
}

func stringToUint(s string) (uint, error) {
	ui, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, apperrors.ErrBadParam
	}
	return uint(ui), nil
}

func parseQueryParamInt(g *gin.Context, queryParamName string) (int, error) {
	str := g.Query(queryParamName)
	if str == "" {
		return 0, apperrors.ErrBadParam
	}
	num, err := strconv.ParseInt(str, 10, 62)
	if err != nil {
		return 0, apperrors.ErrBadParam
	}
	return int(num), nil
}
