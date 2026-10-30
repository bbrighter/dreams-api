package controller

import (
	"net/http"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func handleError(g *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	status := apperrors.HttpStatus(err)

	response := ErrorResponse{
		Code:    http.StatusText(status),
		Message: err.Error(),
	}

	g.AbortWithStatusJSON(status, response)
	return true
}
