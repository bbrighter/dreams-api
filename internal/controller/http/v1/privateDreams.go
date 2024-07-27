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
		h.GET("/", r.GetAll)
		h.GET("/:id", r.Get)
	}
}

// @Description Get all dreams - inlcuding private
// @Produce json
// @Success 200 {object} entity.DreamsResponse "List of all dreams"
// @Router /dreams/private [get]
func (r *privateDreamsRoute) GetAll(g *gin.Context) {
	dreams := r.d.GetAll(true)
	g.JSON(200, dreams.ToResponse())
}

// @Description Get one private dream
// @Produce json
// @Success 200 {object} entity.DreamResponse "One dream"
// @Failure 400
// @Failure 404
// @Router /dreams/private/{dreamId} [get]
func (r *privateDreamsRoute) Get(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	visible := true
	dream, err := r.d.Get(id, visible)
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
