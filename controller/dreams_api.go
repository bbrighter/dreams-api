package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (con Controller) GetDreams(g *gin.Context) {
	dreams := con.Repo.GetDreams()
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
	id, err := con.Repo.CreateDream(body.dreamRequestBodyToDream())
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

	if err := con.Repo.UpdateDream(dream); err != nil {
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
	dream, err := con.Repo.GetDream(uint(id))
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
	err = con.Repo.DeleteDream(uint(id))
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
