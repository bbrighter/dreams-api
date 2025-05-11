package v1

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testDate = time.Date(2020, 11, 30, 13, 45, 52, 0, time.UTC)

type mockDreamsUseCaseOk struct{}

func (tuc mockDreamsUseCaseOk) List() entity.Dreams {
	resp := entity.Dreams{entity.Dream{ID: 1, Date: testDate, Description: "Description", Visible: true}}
	return resp
}

func (tuc mockDreamsUseCaseOk) Get(id uint) (entity.Dream, error) {
	if id == uint(1) {
		return entity.Dream{
			ID:          1,
			Date:        time.Now(),
			Description: "desc",
			Visible:     true,
		}, nil
	}
	return entity.Dream{}, entity.ErrorNotFound
}

func (tuc mockDreamsUseCaseOk) Create(date time.Time) (uint, error) {
	return 1, nil
}

func (tuc mockDreamsUseCaseOk) Update(id uint, date time.Time, description string) error {
	return nil
}

func (tuc mockDreamsUseCaseOk) Delete(id uint) (entity.Categories, error) {
	return entity.Categories{}, nil
}

type mockCategoriesAdderRemover struct{}

func (mu mockCategoriesAdderRemover) AddToDream(categoryName string, dreamId uint) (entity.Categories, error) {
	return entity.Categories{}, nil
}

func (mu mockCategoriesAdderRemover) RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error) {
	return entity.Categories{}, nil
}

type mockPersonsAdderRemover struct{}

func (tu mockPersonsAdderRemover) AddToDream(categoryName string, dreamId uint) (entity.Persons, error) {
	return entity.Persons{}, nil
}
func (tu mockPersonsAdderRemover) RemoveFromDream(categoryId uint, dreamId uint) (entity.Persons, error) {
	return entity.Persons{}, nil
}

func newTestRoute() (*dreamsRoutes, *gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = new(http.Request)
	c.Request.URL = new(url.URL)
	return &dreamsRoutes{
		d: mockDreamsUseCaseOk{},
		c: mockCategoriesAdderRemover{},
		p: mockPersonsAdderRemover{},
	}, c, rec
}

func TestList(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
		uc           usecase.Dreams
	}{
		{
			name:         "ok",
			expectedCode: 200,
			uc:           mockDreamsUseCaseOk{},
			expectedBody: `{"dreams":[{"id":1,"date":"2020-11-30T13:45:52Z","visible":true}]}`,
		},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		route := &dreamsRoutes{d: test.uc, p: mockPersonsAdderRemover{}, c: mockCategoriesAdderRemover{}}

		t.Run(test.name, func(t *testing.T) {
			route.GetAll(c)

			assert.Equal(t, test.expectedCode, rec.Code)
			assert.Equal(t, test.expectedBody, rec.Body.String())
		})
	}

	// 	r, g, rec := newTestRoute()
	// 	r.GetAll(g)

	// 	assert.Equal(t, rec.Code, 200)
	// 	var response map[string][]string
	// 	err := json.Unmarshal(rec.Body.Bytes(), &response)
	// 	assert.NoError(t, err)
	// 	value, exists := response["dreams"]
	// 	assert.True(t, exists)
	// 	assert.Len(t, value, 0)
}

func TestGet(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")
	r.Get(g)

	assert.Equal(t, 200, rec.Code)
}

func TestGet404(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "2")

	r.Get(g)
	assert.Equal(t, 404, rec.Code)
}

func TestGet400(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "aaa")

	r.Get(g)
	assert.Equal(t, 400, rec.Code)
}

func TestCreate(t *testing.T) {
	// Don't know how to add a body yet
	t.Skip()
	r, g, rec := newTestRoute()
	var body = strings.NewReader(`"description":"Description","date":"2022-02-1"`)
	g.Request = new(http.Request)
	// req, _ := http.NewRequest("POST", "/", body)
	g.Request.Body = io.NopCloser(body)
	// g.Request = req
	r.Create(g)

	assert.Equal(t, 200, rec.Code)
}

func TestDelete(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")

	r.Delete(g)

	assert.Equal(t, 200, rec.Code)
}

func TestAddCategoryToDream(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")
	g.Request.URL, _ = url.Parse("?name=Name")
	r.PutCategoryToDream(g)

	assert.Equal(t, 200, rec.Code)
}

func TestAddCategoryToDream400NoName(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")
	r.PutCategoryToDream(g)

	assert.Equal(t, 400, rec.Code)
}

func TestAddCategoryToDream400NoID(t *testing.T) {
	r, g, rec := newTestRoute()
	g.Request.URL, _ = url.Parse("?name=Name")
	r.PutCategoryToDream(g)

	assert.Equal(t, 400, rec.Code)
}

func TestRemoveCategoryFromDream(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")
	g.AddParam("categoryId", "1")
	r.RemoveCategoryFromDream(g)

	assert.Equal(t, 200, rec.Code)
}

func TestRemoveCategoryFromDream400NoID(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("categoryId", "1")
	r.RemoveCategoryFromDream(g)

	assert.Equal(t, 400, rec.Code)
}
func TestRemoveCategoryFromDream400NoCategoryID(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")
	r.RemoveCategoryFromDream(g)

	assert.Equal(t, 400, rec.Code)
}

func TestAddPersonToDream(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")
	g.Request.URL, _ = url.Parse("?name=Name")
	r.PutPersonToDream(g)

	assert.Equal(t, 200, rec.Code)
}

func TestAddPersonToDream400NoID(t *testing.T) {
	r, g, rec := newTestRoute()
	g.Request.URL, _ = url.Parse("?name=Name")
	r.PutPersonToDream(g)

	assert.Equal(t, 400, rec.Code)
}

func TestAddPersonToDream400NoName(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")
	r.PutPersonToDream(g)

	assert.Equal(t, 400, rec.Code)
}

func TestRemovePersonFromDream(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")
	g.AddParam("personId", "1")
	r.RemovePersonFromDream(g)

	assert.Equal(t, 200, rec.Code)
}

func TestRemovePersonFromDream400NoID(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("personId", "1")
	r.RemovePersonFromDream(g)

	assert.Equal(t, 400, rec.Code)
}

func TestRemovePersonFromDream400NoPersonID(t *testing.T) {
	r, g, rec := newTestRoute()
	g.AddParam("id", "1")

	r.RemovePersonFromDream(g)

	assert.Equal(t, 400, rec.Code)
}
