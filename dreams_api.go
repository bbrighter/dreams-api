package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DreamRequestBody struct {
	Date        time.Time `json:"date"`
	Description *string   `json:"description"`
}

func (body DreamRequestBody) dreamRequestBodyToDream() Dream {
	var desc string
	if body.Description != nil {
		desc = *body.Description
	}
	return Dream{
		Date:        body.Date,
		Description: desc,
	}
}

type DreamResponse struct {
	ID          uint      `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

type DreamsResponse struct {
	Dreams []DreamResponse `json:"dreams"`
}

func dreamToDreamResponse(d Dream) DreamResponse {
	return DreamResponse{
		ID:          d.ID,
		Date:        d.Date,
		Description: d.Description,
	}
}

func dreamsToDreamsResponse(dreams []Dream) DreamsResponse {
	var respList = []DreamResponse{}
	for _, d := range dreams {
		dream := dreamToDreamResponse(d)
		respList = append(respList, dream)
	}
	return DreamsResponse{Dreams: respList}
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
	if body.Date.IsZero() {
		g.AbortWithError(http.StatusBadRequest, ErrorParameterMissing("date"))
		return
	}
	id, err := con.Repo.createDream(body.dreamRequestBodyToDream())
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.JSON(http.StatusCreated, id)
}

func (con Controller) UpdateDream(g *gin.Context) {
	var body DreamRequestBody
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var dream = body.dreamRequestBodyToDream()

	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dream.ID = uint(id)

	if err := con.Repo.updateDream(dream); err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.Status(http.StatusOK)
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
