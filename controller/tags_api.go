package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagsResponse = []TagResponse

type TagResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func (con Controller) GetTags(g *gin.Context) {
	tags := tagsToTagsResponse(con.Repo.GetTags())
	g.JSON(http.StatusOK, tags)
}

func (con Controller) AddTag(g *gin.Context) {
	dreamId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	tagTitle := g.Query("title")
	if tagTitle == "" {
		g.AbortWithError(http.StatusBadRequest, ErrorParameterMissing("title"))
		return
	}

	id, err := con.Repo.AddTagToDream(tagTitle, uint(dreamId))
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, id)
}

func (con Controller) RemoveTag(g *gin.Context) {
	dreamId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	tagId, err := strconv.Atoi(g.Param("tagId"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tags, err := con.Repo.RemoveTagFromDream(uint(tagId), uint(dreamId))
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, tagsToTagsResponse(tags))
}
