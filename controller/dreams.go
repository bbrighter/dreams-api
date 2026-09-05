package controller

import (
	"net/http"
	"time"

	dreamcategories "github.com/bbrighter/dreams-api/internal/dreamCategories"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/entities"

	"github.com/gin-gonic/gin"
)

type dreamsRoutes struct {
	d *dreams.DreamsService
	c *dreamcategories.DreamCategoriesService
}

func newDreamsRoute(
	handler *gin.RouterGroup,
	d *dreams.DreamsService,
	c *dreamcategories.DreamCategoriesService,
) {
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

type DreamListResponse struct {
	Dreams []DreamResponse `json:"dreams"`
}
type DreamResponse struct {
	ID          uint               `json:"id"`
	Date        time.Time          `json:"date"`
	Description string             `json:"description"`
	Finalized   bool               `json:"finalized"`
	Categories  []CategoryResponse `json:"categories"`
	Rating      *int               `json:"rating,omitempty"`
}

func dreamToResponse(dream entities.Dream) DreamResponse {
	return DreamResponse{
		ID:          dream.ID,
		Date:        dream.Date,
		Description: dream.Description,
		Finalized:   dream.Finalized,
		Categories:  categoriesToResponse(dream.Categories).Categories,
		Rating:      dream.Rating,
	}
}

func dreamsToResponse(dreams []entities.Dream) DreamListResponse {
	var resp = []DreamResponse{}
	for _, dream := range dreams {
		resp = append(resp, dreamToResponse(dream))
	}
	return DreamListResponse{Dreams: resp}
}

// @Description Get all dreams
// @Produce json
// @Success 200 {object} DreamListResponse "List of all dreams"
// @Router /dreams [get]
func (r *dreamsRoutes) GetAll(g *gin.Context) {
	dreams, err := r.d.ListDreams(g.Request.Context())
	if handleError(g, err) {
		return
	}
	g.JSON(200, dreamsToResponse(dreams))
}

// @Description Get one dream
// @Produce json
// @Success 200 {object} DreamResponse "One dream"
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId} [get]
func (r *dreamsRoutes) Get(g *gin.Context) {
	id, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}

	dream, err := r.d.GetDreamById(g.Request.Context(), id)
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, dreamToResponse(dream))
}

type PostDreamRequest struct {
	Date        time.Time `json:"date"`
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
		g.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := r.d.CreateDream(g.Request.Context(), body.Date)
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

	err = r.d.UpdateDream(g.Request.Context(), id, dreams.UpdateDreamParams{
		Date:        body.Date,
		Description: body.Description,
		Rating:      body.Rating,
	})
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

	err = r.d.DeleteDream(g.Request.Context(), id)
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
// @Router /dreams/{dreamId}/categories/{categoryId} [put]
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
	Name string `json:"name"`
	Type string `json:"categoryType"`
}

// @Description Add a new person or category to a dream
// @Produce json
// @Success 201 {number} ID
// @Failure 400
// @Failure 404
// @Router /dreams/{dreamId}/categories [post]
// @Param params body PostCategoryRequestBody true "Name and category. Allowed values for 'category' are: 'person', 'category'"
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
	catType := entities.CategoryType(body.Type)
	id, err := r.c.AddNewCategoryToDream(g.Request.Context(), dreamId, body.Name, catType)
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
// @Router /dreams/{dreamId}/categories/{categoryId} [delete]
func (r *dreamsRoutes) RemoveCategoryFromDream(g *gin.Context) {
	dreamId, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}
	catId, err := parseParamUint(g, "categoryId")
	if handleError(g, err) {
		return
	}

	err = r.c.RemoveCategoryToDream(g.Request.Context(), dreamId, catId)
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
	var finalized = true

	err = r.d.UpdateDream(g.Request.Context(), dreamId, dreams.UpdateDreamParams{Finalized: &finalized})
	if handleError(g, err) {
		return
	}

	g.Status(http.StatusOK)
}
