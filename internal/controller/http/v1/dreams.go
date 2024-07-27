package v1

import (
	"net/http"
	"time"

	customerrors "github.com/bbrighter/dreams-api/internal/customErrors"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type dreamsRoutes struct {
	d usecase.Dreams
	p usecase.Persons
	c usecase.Categories
}

func newDreamsRoute(handler *gin.RouterGroup, d usecase.Dreams, p usecase.Persons, c usecase.Categories) {
	r := &dreamsRoutes{d, p, c}

	h := handler.Group("/dreams")
	{
		h.GET("/", r.GetAll)
		h.POST("/", r.Create)
	}

	hid := h.Group("/:id")
	{
		hid.GET("/", r.Get)
		hid.PATCH("/", r.Update)
		hid.DELETE("/", r.Delete)
	}

	hcat := hid.Group("/category")
	{
		hcat.PUT("/", r.PutCategoryToDream)
		hcat.DELETE("/:categoryId", r.RemoveCategoryFromDream)
	}

	hper := hid.Group("/person")
	{
		hper.PUT("/", r.PutPersonToDream)
		hper.DELETE("/:personId", r.RemovePersonFromDream)
	}
}

// @Description Get all dreams
// @Produce json
// @Success 200 {object} entity.DreamsResponse "List of all dreams"
// @Router /dreams [get]
func (r *dreamsRoutes) GetAll(g *gin.Context) {
	dreams := r.d.GetAll(false)
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
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dream, err := r.d.Get(id, false)
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
		g.AbortWithError(http.StatusBadRequest, customerrors.ErrorParameterMissing("date"))
		return
	}
	id, err := r.d.Create(body.Date)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
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
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = r.d.Update(id, body.Date, *body.Description)
	if err != nil && err == customerrors.ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
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
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	categories, err := r.d.Delete(id)
	if err != nil && err == customerrors.ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, categories.ToResponse())
}

// @Description Add a person to a dream
// @Produce json
// @Success 200 {number} entity.PersonsResonse
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/persons [put]
// @Param name query string true "Name of person"
func (r *dreamsRoutes) PutPersonToDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "id")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	name, err := parseQueryParamString(g, "name")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	persons, err := r.p.AddToDream(name, dreamId)
	if err == customerrors.ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
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
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	personId, err := parseParamUint(g, "personId")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	persons, err := r.p.RemoveFromDream(personId, dreamId)
	if err == customerrors.ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
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
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	name, err := parseQueryParamString(g, "name")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	categories, err := r.c.AddToDream(name, dreamId)
	if err == customerrors.ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
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
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	categoryId, err := parseParamUint(g, "categoryId")
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	categories, err := r.c.RemoveFromDream(categoryId, dreamId)
	if err == customerrors.ErrorNotFound {
		g.AbortWithStatus(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, categories.ToResponse())
}
