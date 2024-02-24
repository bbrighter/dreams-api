package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DreamMetaResponse struct {
	ID      uint      `json:"id" validate:"required"`
	Date    time.Time `json:"date" validate:"required"`
	Visible bool      `json:"visible" validate:"required"`
}

// Respone when querying one dream
type DreamResponse struct {
	DreamMetaResponse
	Description string             `json:"description" validate:"required"`
	Categories  []CategoryResponse `json:"categories" validate:"required"`
	Persons     []PersonResponse   `json:"persons" validate:"required"`
}

// Returned when querying all dreams
type DreamsResponse struct {
	Dreams []DreamMetaResponse `json:"dreams" validate:"required"`
}

// @Description Get all dreams
// @Produce json
// @Success 200 {object} DreamsResponse "List of all dreams"
// @Param showPrivateDreams query bool false "True if all dreams should be shown"
// @Router /dreams [get]
func GetDreams(g *gin.Context) {
	con := GetCon(g)
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
// @Success 201 {number} ID
// @Failure 400
// @Router /dreams [post]
// @Param dreamRequestBody  body DreamRequestBody true "The dream which will be created"
func CreateDream(g *gin.Context) {
	con := GetCon(g)
	var body DreamRequestBody
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if body.Date.IsZero() {
		g.AbortWithError(http.StatusBadRequest, ErrorParameterMissing("date"))
		return
	}
	id, _ := con.Repo.CreateDream(body.dreamRequestBodyToDream())
	g.JSON(http.StatusCreated, id)
}

// @Description Update an existing dream
// @Accept json
// @Success 200
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [patch]
// @Param dreamRequestBody body DreamRequestBody true "The dream which will be updated"
func UpdateDream(g *gin.Context) {
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

	con := GetCon(g)
	err = con.Repo.UpdateDream(dream)
	if err != nil && err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.Status(http.StatusOK)
}

// @Description Get one dream
// @Produce json
// @Success 200 {object} DreamResponse "One dream"
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [get]
func GetDream(g *gin.Context) {
	getDreamOrNot(g, false)
}

// @Description Delete one dreams
// @Produce json
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [delete]
func DeleteDream(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	con := GetCon(g)
	err = con.Repo.DeleteDream(id)
	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.Status(http.StatusOK)
}

// @Description Get all private dreams
// @Produce json
// @Success 200 {object} DreamsResponse "List of all dreams"
// @Failure 401
// @Security BasicAuth
// @Param Authorization header string true "Basic Username:password"
// @Router /private/dreams [get]
func GetPrivateDreams(g *gin.Context) {
	var showAll = true
	con := GetCon(g)
	dreams := con.Repo.GetDreams(&showAll)
	g.JSON(http.StatusOK, dreamsToDreamsResponse(dreams))
}

// @Description Get one private dream
// @Produce json
// @Success 200 {object} DreamResponse "One dream"
// @Failure 400
// @Failure 401
// @Failure 404
// @Security BasicAuth
// @Router /private/dreams/{dreamId} [get]
func GetPrivateDream(g *gin.Context) {
	getDreamOrNot(g, true)
}

// @Description Toggle visiblity of a dream
// @Produce json
// @Success 200 {bool} isPrivate
// @Failure 400
// @Failure 401
// @Failure 404
// @Security BasicAuth
// @Router /private/dreams/{dreamId} [patch]
func TogglePrivateDream(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	con := GetCon(g)
	visible, err := con.Repo.TogglePrivateDream(id)
	if err != nil && err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, visible)
}

func getDreamOrNot(g *gin.Context, showAll bool) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	con := GetCon(g)
	dream, err := con.Repo.GetDream(id, &showAll)

	if err == ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, dreamToDreamResponse(dream))
}
