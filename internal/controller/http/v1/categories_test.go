package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockCategoriesLister struct{}

func (tu mockCategoriesLister) List() entity.Categories {
	return entity.Categories{entity.Category{
		ID: 1, Name: "Name", Dreams: []entity.Dream{{ID: 10}},
	}}
}

func TestGetAllCategories(t *testing.T) {
	tests := []struct {
		name           string
		uc             usecase.CategoriesLister
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Ok",
			uc:             mockCategoriesLister{},
			expectedStatus: 200,
			expectedBody:   `{"categories":[{"id":1,"name":"Name"}]}`,
		},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		route := &categoriesRoute{c: test.uc}

		t.Run(test.name, func(t *testing.T) {
			route.GetAll(c)

			assert.Equal(t, rec.Code, test.expectedStatus)
			assert.Equal(t, rec.Body.String(), test.expectedBody)
		})
	}
}
