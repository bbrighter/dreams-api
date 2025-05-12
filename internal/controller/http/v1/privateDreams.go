package v1

import (
	"net/http"

	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type privateDreamsRoute struct {
	p usecase.PrivateDreams
}

func newPrivateDreamsRoute(handler *gin.RouterGroup, p usecase.PrivateDreams) {
	r := &privateDreamsRoute{p: p}

	h := handler.Group("/dreams/private")
	{
		h.GET("", r.GetAll)
		h.GET("/:id", r.Get)
		h.PATCH("/:id", r.ToggleVisibility)
	}
}

// @Description Get all dreams - including private
// @Produce json
// @Success 200 {object} entity.DreamsResponse "List of all dreams"
// @Router /dreams-api/v1/dreams/private [get]
func (r *privateDreamsRoute) GetAll(g *gin.Context) {
	dreams := r.p.List()
	g.JSON(200, dreams.ToResponse())
}

// @Description Get one private dream
// @Produce json
// @Success 200 {object} entity.DreamResponse "One dream"
// @Failure 400
// @Failure 404
// @Router /dreams-api/v1/dreams/private/{dreamId} [get]
func (r *privateDreamsRoute) Get(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		return
	}
	dream, err := r.p.Get(id)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, dream.ToResponse())
}

// @Description Toggle visibility of a dream
// @Produce json
// @Success 200
// @Router /dreams-api/v1/dreams/private/{dreamId} [patch]
func (r *privateDreamsRoute) ToggleVisibility(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		return
	}
	err = r.p.ToggleVisibility(id)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
}
