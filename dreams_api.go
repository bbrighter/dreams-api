package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DreamRequestBody struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

func (body DreamRequestBody) dreamRequestBodyToDream() Dream {
	return Dream{
		Date:        body.Date,
		Description: body.Description,
	}
}

type DreamResponse struct {
	ID          uint      `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

func dreamToDreamResponse(d Dream) DreamResponse {
	return DreamResponse{
		ID:          d.ID,
		Date:        d.Date,
		Description: d.Description,
	}
}

func dreamsToDreamsResponse(dreams []Dream) []DreamResponse {
	var dreamsResp []DreamResponse
	for _, d := range dreams {
		dream := dreamToDreamResponse(d)
		dreamsResp = append(dreamsResp, dream)
	}
	return dreamsResp
}

func (con Controller) GetDreams(g *gin.Context) {
	dreams, err := con.Repo.getDreams()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.JSON(http.StatusOK, dreamsToDreamsResponse(dreams))
}

func (con Controller) CreateDream(g *gin.Context) {
	var body DreamRequestBody
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id, err := con.Repo.createDream(body.dreamRequestBodyToDream())
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.JSON(http.StatusCreated, id)
}

func (con Controller) GetDream(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dream, err := con.Repo.getDream(uint(id))
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

func (con Controller) DeleteDream(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = con.Repo.deleteDream(uint(id))
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
