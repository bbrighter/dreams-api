package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testUseCasePersons struct{}

func (tu testUseCasePersons) GetAll() entity.Persons {
	return entity.Persons{entity.Person{
		ID: 1, Name: "Name", Dreams: []entity.Dream{{ID: 10}},
	}}
}

func (tu testUseCasePersons) AddToDream(categoryName string, dreamId uint) (entity.Persons, error) {
	return entity.Persons{}, nil
}

func (tu testUseCasePersons) RemoveFromDream(categoryId uint, dreamId uint) (entity.Persons, error) {
	return entity.Persons{}, nil
}

func newTestRoutePerson() (*personsRoute, *gin.Context, *httptest.ResponseRecorder) {
	var tuc = testUseCasePersons{}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	return &personsRoute{p: tuc}, c, rec
}

func TestGetAllPersons(t *testing.T) {
	r, g, rec := newTestRoutePerson()

	r.GetAll(g)

	assert.Equal(t, rec.Code, 200)
	assert.Equal(t, rec.Body.String(), `{"persons":[{"id":1,"name":"Name"}]}`)
}
