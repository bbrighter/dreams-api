package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DreamMetaResponse struct {
	ID   uint      `json:"id" validate:"required"`
	Date time.Time `json:"date" validate:"required"`
}

type DreamResponse struct {
	DreamMetaResponse
	Description string             `json:"description" validate:"required"`
	Categories  []CategoryResponse `json:"categories" validate:"required"`
	Persons     []PersonResponse   `json:"persons" validate:"required"`
}

type DreamsResponse struct {
	Dreams []DreamMetaResponse `json:"dreams" validate:"required"`
}

// @Description Get all dreams
// @Produce json
// @Success 200 {object} DreamsResponse "List of all dreams"
// @Router /dreams [get]
func (con Controller) GetDreams(g *gin.Context) {
	var showAll bool = false
	dreams := con.Repo.GetDreams(&showAll)
	g.JSON(http.StatusOK, dreamsToDreamsResponse(dreams))
}

type DreamRequestBody struct {
	Date        time.Time `json:"date" validate:"required"`
	Description *string   `json:"description"`
}

// @Description Create a new dream
// @Accept json
// @Produce json
// @Success 200 {number} ID
// @Router /dreams [post]
// @Param dreamRequestBody  body DreamRequestBody true "The dream which will be created"
func (con Controller) CreateDream(g *gin.Context) {
	var body DreamRequestBody
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if body.Date.IsZero() {
		g.AbortWithError(http.StatusBadRequest, ErrorParameterMissing("date"))
		return
	}
	id, err := con.Repo.CreateDream(body.dreamRequestBodyToDream())
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.JSON(http.StatusCreated, id)
}

// @Description Update an existing dream
// @Accept json
// @Success 200
// @Router /dreams/{dreamId} [patch]
// @Param dreamRequestBody body DreamRequestBody true "The dream which will be updated"
func (con Controller) UpdateDream(g *gin.Context) {
	var body DreamRequestBody
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var dream = body.dreamRequestBodyToDream()

	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dream.ID = id

	if err := con.Repo.UpdateDream(dream); err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.Status(http.StatusOK)
}

// @Description Get one dream
// @Produce json
// @Success 200 {object} DreamResponse "One dream"
// @Router /dreams/{dreamId} [get]
func (con Controller) GetDream(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var showAll bool = false
	dream, err := con.Repo.GetDream(id, &showAll)

	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.JSON(http.StatusOK, dreamToDreamResponse(dream))
}

// @Description Delete one dreams
// @Produce json
// @Success 200
// @Router /dreams/{dreamId} [delete]
func (con Controller) DeleteDream(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = con.Repo.DeleteDream(id)
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.Status(http.StatusOK)
}
