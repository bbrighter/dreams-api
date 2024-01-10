package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagsResponse struct {
	Tags []TagResponse `json:"tags" validate:"required"`
}

type TagResponse struct {
	ID    uint   `json:"id" validate:"required"`
	Title string `json:"title" validate:"required"`
}

// @Description Get all tags
// @Produce json
// @Success 200 {object} TagsResponse
// @Router /tags [get]
func (con Controller) GetTags(g *gin.Context) {
	tags := tagsToTagsResponse(con.Repo.GetTags())
	g.JSON(http.StatusOK, tags)
}

// @Description Add a tag to a dream
// @Produce json
// @Success 200 {object} TagsResponse
// @Router /dream/{dreamId}/tags [put]
// @Param title query number true "Label of tag"
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

	tags, err := con.Repo.AddTagToDream(tagTitle, uint(dreamId))
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, tagsToTagsResponse(tags))
}

// @Description Remove a tag to a dream
// @Produce json
// @Success 200 {object} TagsResponse
// @Router /dream/{dreamId}/tags/{tagId} [delete]
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
