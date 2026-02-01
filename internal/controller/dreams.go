package controller

import (
	"net/http"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type dreamsRoutes struct {
	d usecase.IDreams
	c usecase.ICategories
}

func newDreamsRoute(handler *gin.RouterGroup, d usecase.IDreams, c usecase.ICategories) {
	r := &dreamsRoutes{d, c}

	h := handler.Group("/dreams")
	{
		h.GET("", r.GetAll)
		h.POST("", r.Create)
		hid := h.Group("/:dreamId")
		{
			hid.GET("", r.Get)
			hid.PATCH("", r.Update)
			hid.DELETE("", r.Delete)
			hid.PATCH("/finalize", r.Finalize)
			hCat := hid.Group("/categories")
			{
				hCat.POST("", r.AddCategoryToDream)
				hCat.PUT("/:categoryId", r.PutCategoryToDream)
				hCat.DELETE("/:categoryId", r.RemoveCategoryFromDream)
			}
		}
	}
}

// @Description Get all dreams
// @Produce json
// @Success 200 {object} entity.DreamsResponse "List of all dreams"
// @Router /dreams [get]
func (r *dreamsRoutes) GetAll(g *gin.Context) {
	dreams, _ := r.d.List(g.Request.Context())
	g.JSON(200, dreams.ToResponse())
}

// @Description Get one dream
// @Produce json
// @Success 200 {object} entity.DreamResponse "One dream"
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [get]
func (r *dreamsRoutes) Get(g *gin.Context) {
	id, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}

	dream, err := r.d.Get(g.Request.Context(), id)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, dream.ToResponse())
}

type PostDreamRequest struct {
	Date        time.Time `json:"date" binding:"required"`
	Description *string   `json:"description"`
}

// @Description Create a new dream
// @Accept json
// @Produce json
// @Success 201 {number} ID
// @Failure 400
// @Failure 500
// @Router /dreams [post]
// @Param postDreamRequest  body PostDreamRequest true "The dream which will be created"
func (r *dreamsRoutes) Create(g *gin.Context) {
	var body PostDreamRequest
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if body.Date.IsZero() {
		g.AbortWithError(http.StatusBadRequest, entity.ErrorBadParam)
		return
	}
	id, err := r.d.Create(g.Request.Context(), body.Date)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusCreated, id)
}

type UpdateDreamRequest struct {
	Date        *time.Time `json:"date,omitempty"`
	Description *string    `json:"description,omitempty"`
	Rating      *int       `json:"rating,omitempty"`
}

// @Description Update an existing dream
// @Accept json
// @Success 200
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [patch]
// @Param updateDreamRequest body UpdateDreamRequest true "All parameters of the dream that should be updated"
func (r *dreamsRoutes) Update(g *gin.Context) {
	var body UpdateDreamRequest
	err := g.BindJSON(&body)
	if handleError(g, err) {
		return
	}
	id, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}

	err = r.d.Update(g.Request.Context(), id, body.Date, body.Description, body.Rating)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
}

// @Description Delete one dreams
// @Produce json
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [delete]
func (r *dreamsRoutes) Delete(g *gin.Context) {
	id, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}

	err = r.d.Delete(g.Request.Context(), id)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
}

// @Description Add an exiting person or category to a dream
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/category/{categoryId} [put]
func (r *dreamsRoutes) PutCategoryToDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}
	catId, err := parseParamUint(g, "categoryId")
	if handleError(g, err) {
		return
	}
	err = r.c.AddCategoryToDream(g.Request.Context(), dreamId, catId)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
}

type PostCategoryRequestBody struct {
	Name string              `json:"name"`
	Type entity.CategoryType `json:"categoryType"`
}

// @Description Add a new person or category to a dream
// @Produce json
// @Success 201 {number} ID
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/category [post]
// @Param params body PostCategoryRequestBody true "Name and category. Allowed values for 'catgory' are: 'person', 'category'"
func (r *dreamsRoutes) AddCategoryToDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}
	var body PostCategoryRequestBody
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id, err := r.c.AddNewCategoryToDream(g.Request.Context(), dreamId, body.Name, body.Type)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusCreated, id)

}

// @Description Delete a category or person from a dream
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/category/{categoryId} [delete]
func (r *dreamsRoutes) RemoveCategoryFromDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}
	catId, err := parseParamUint(g, "categoryId")
	if handleError(g, err) {
		return
	}

	err = r.c.RemoveCategoryFromDream(g.Request.Context(), dreamId, catId)
	if handleError(g, err) {
		return
	}
	g.Status(http.StatusOK)
}

// @Description Finalize a dream
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/finalize [patch]
func (r *dreamsRoutes) Finalize(g *gin.Context) {
	dreamId, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}

	err = r.d.Finalize(g.Request.Context(), dreamId)
	if handleError(g, err) {
		return
	}

	g.Status(http.StatusOK)
}
