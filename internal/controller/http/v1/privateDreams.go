package v1

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type privateDreamsRoute struct {
	d usecase.Dreams
}

func newPrivateDreamsRoute(handler *gin.RouterGroup, d usecase.Dreams) {
	r := &privateDreamsRoute{d}

	h := handler.Group("/dreams/private")
	{
		h.GET("", r.GetAll)
		h.GET("/:id", r.Get)
		h.PATCH("/:id", r.ToggleVisibility)
	}
}

// @Description Get all dreams - inlcuding private
// @Produce json
// @Success 200 {object} entity.DreamsResponse "List of all dreams"
// @Router /v1/dreams/private [get]
func (r *privateDreamsRoute) GetAll(g *gin.Context) {
	dreams := r.d.GetAll(true)
	g.JSON(200, dreams.ToResponse())
}

// @Description Get one private dream
// @Produce json
// @Success 200 {object} entity.DreamResponse "One dream"
// @Failure 400
// @Failure 404
// @Router /v1/dreams/private/{dreamId} [get]
func (r *privateDreamsRoute) Get(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dream, err := r.d.Get(id, true)
	if err != nil && err.Error() == "record not found" {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.JSON(http.StatusOK, dream.ToResponse())
}

// @Description Toggle visiblity of a dream
// @Produce json
// @Success 200
// @Router /v1/dreams/private/{dreamId} [patch]
func (r *privateDreamsRoute) ToggleVisibility(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err = r.d.ToggleVisibility(id); err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.Status(http.StatusOK)
}
