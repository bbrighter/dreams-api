package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryCount struct {
	CategoryId uint `json:"categoryId" validate:"required"`
	Count      int  `json:"count" validate:"required"`
}

type CategoryCountResponse struct {
	Categories []CategoryCount `json:"categories" validate:"required"`
}

// @Description Get count per category
// @Produce json
// @Success 200 {object} CategoryCountResponse "Count by categoryId"
// @Router /statistics [get]
func (con Controller) GetCountCategories(g *gin.Context) {
	count := con.Repo.CountCategories()
	g.JSON(http.StatusOK, categoryCountsToCategoryCountResponse(count))
}
