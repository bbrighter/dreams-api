package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testUseCaseCategories struct{}

func (tu testUseCaseCategories) GetAll() entity.Categories {
	return entity.Categories{entity.Category{
		ID: 1, Name: "Name", Dreams: []entity.Dream{{ID: 10}},
	}}
}

func (tu testUseCaseCategories) AddToDream(categoryName string, dreamId uint) (entity.Categories, error) {
	return entity.Categories{}, nil
}

func (tu testUseCaseCategories) RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error) {
	return entity.Categories{}, nil
}

func newTestRouteCat() (*categoriesRoute, *gin.Context, *httptest.ResponseRecorder) {
	var tuc = testUseCaseCategories{}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	return &categoriesRoute{c: tuc}, c, rec
}

func TestGetAllCategories(t *testing.T) {
	r, g, rec := newTestRouteCat()
	r.GetAll(g)

	assert.Equal(t, rec.Code, 200)
	assert.Equal(t, rec.Body.String(), `{"categories":[{"id":1,"name":"Name"}]}`)
}
