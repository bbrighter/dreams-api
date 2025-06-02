package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

func TestCategoriesAndPersons(t *testing.T) {
	g := setupApiTest(t)
	date := time.Now()

	tests := []apiTest{
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: DreamRequestBody{Date: date}},
		{name: "add category to dream", method: http.MethodPut, url: "/dreams/1/categories?name=cat", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 1, Name: "cat"}}}},
		{name: "get categories", method: http.MethodGet, url: "/categories", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 1, Name: "cat"}}}},
		{name: "remove category from dream", method: http.MethodDelete, url: "/dreams/1/categories/1", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{}}},
		{name: "get categories", method: http.MethodGet, url: "/categories", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{}}},
	}
	for _, test := range tests {
		test.evaluate(t, g)
	}
}
