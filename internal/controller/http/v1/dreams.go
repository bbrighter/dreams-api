package v1

import (
	"net/http"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type dreamsRoutes struct {
	d usecase.Dreams
	p usecase.PersonsAdderRemover
	c usecase.CategoriesAdderRemover
}

func newDreamsRoute(handler *gin.RouterGroup, d usecase.Dreams, p usecase.PersonsAdderRemover, c usecase.CategoriesAdderRemover) {
	r := &dreamsRoutes{d, p, c}

	h := handler.Group("/dreams")
	{
		h.GET("", r.GetAll)
		h.POST("", r.Create)
		hid := h.Group("/:id")
		{
			hid.GET("", r.Get)
			hid.PATCH("", r.Update)
			hid.DELETE("", r.Delete)
			hid.PATCH("/finalize", r.Finalize)
			hCat := hid.Group("/categories")
			{
				hCat.PUT("", r.PutCategoryToDream)
				hCat.DELETE("/:categoryId", r.RemoveCategoryFromDream)
			}
			hPer := hid.Group("/persons")
			{
				hPer.PUT("", r.PutPersonToDream)
				hPer.DELETE("/:personId", r.RemovePersonFromDream)
			}
		}
	}
}

// @Description Get all dreams
// @Produce json
// @Success 200 {object} entity.DreamsResponse "List of all dreams"
// @Router /dreams [get]
// @Param includes query string false "Comma separated list of child objects. Possible entries: categories, persons"
func (r *dreamsRoutes) GetAll(g *gin.Context) {
	includes := entity.ParseIncludes(g.Query("includes"))
	dreams := r.d.List(includes)
	g.JSON(200, dreams.ToResponse())
}

// @Description Get one dream
// @Produce json
// @Success 200 {object} entity.DreamResponse "One dream"
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [get]
func (r *dreamsRoutes) Get(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if err != nil {
		return
	}

	dream, err := r.d.Get(id)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, dream.ToResponse())
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
// @Failure 500
// @Router /dreams [post]
// @Param dreamRequestBody  body DreamRequestBody true "The dream which will be created"
func (r *dreamsRoutes) Create(g *gin.Context) {
	var body DreamRequestBody
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if body.Date.IsZero() {
		g.AbortWithError(http.StatusBadRequest, entity.ErrorBadParam)
		return
	}
	id, err := r.d.Create(body.Date)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusCreated, id)
}

// @Description Update an existing dream
// @Accept json
// @Success 200
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [patch]
// @Param dreamRequestBody body DreamRequestBody true "The dream which will be updated"
func (r *dreamsRoutes) Update(g *gin.Context) {
	var body DreamRequestBody
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var err error
	var id uint
	id, err = parseParamUint(g, "id")
	if err != nil {
		return
	}
	err = r.d.Update(id, body.Date, *body.Description)
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
	id, err := parseParamUint(g, "id")
	if err != nil {
		return
	}

	categories, _, err := r.d.Delete(id)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, categories.ToResponse())
}

// @Description Add a person to a dream
// @Produce json
// @Success 200 {object} entity.PersonsResponse
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/persons [put]
// @Param name query string true "Name of person"
func (r *dreamsRoutes) PutPersonToDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "id")
	if err != nil {
		return
	}
	name, err := parseQueryParamString(g, "name")
	if err != nil {
		return
	}

	persons, err := r.p.AddToDream(name, dreamId)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, persons.ToResponse())
}

// @Description Delete a person from a dream
// @Produce json
// @Success 200 {object} entity.PersonsResponse
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/persons/{personId} [delete]
func (r *dreamsRoutes) RemovePersonFromDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "id")
	if err != nil {
		return
	}
	personId, err := parseParamUint(g, "personId")
	if err != nil {
		return
	}

	persons, err := r.p.RemoveFromDream(personId, dreamId)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, persons.ToResponse())
}

// @Description Add a category to a dream
// @Produce json
// @Success 200 {object} entity.CategoriesResponse
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/categories [put]
// @Param name query string true "Name of a category"
func (r *dreamsRoutes) PutCategoryToDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "id")
	if err != nil {
		return
	}
	name, err := parseQueryParamString(g, "name")
	if err != nil {
		return
	}

	categories, err := r.c.AddToDream(name, dreamId)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, categories.ToResponse())
}

// @Description Remove a category from a dream
// @Produce json
// @Success 200 {object} entity.CategoriesResponse
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/categories/{categoryId} [delete]
func (r *dreamsRoutes) RemoveCategoryFromDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "id")
	if err != nil {
		return
	}
	categoryId, err := parseParamUint(g, "categoryId")
	if err != nil {
		return
	}

	categories, err := r.c.RemoveFromDream(categoryId, dreamId)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, categories.ToResponse())
}

// @Description Finalize a dream
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/finalize [patch]
func (r *dreamsRoutes) Finalize(g *gin.Context) {
	dreamId, err := parseParamUint(g, "id")
	if err != nil {
		return
	}

	err = r.d.Finalize(dreamId)
	if handleError(g, err) {
		return
	}

	g.Status(http.StatusOK)
}
