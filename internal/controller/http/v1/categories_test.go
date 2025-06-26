package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

func TestCategories(t *testing.T) {
	g := setupApiTest(t)
	date := time.Now()

	var emptyCategories = entity.CategoriesResponse{}

	tests := []apiTest{
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: PostDreamRequest{Date: date}},
		{name: "add category to dream", method: http.MethodPut, url: "/dreams/1/categories?name=cat", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 1, Name: "cat"}}}},
		{name: "change type of category", method: http.MethodPatch, url: "/categories/1?newType=person", statusCode: http.StatusOK},
		{name: "get categories", method: http.MethodGet, url: "/categories", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Persons: []entity.CategoryResponse{{ID: 1, Name: "cat"}}}},
		{name: "get categories and number of dreams", method: http.MethodGet, url: "/categories?includes=dreamsCount", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Persons: []entity.CategoryResponse{{ID: 1, Name: "cat", Count: 1}}}},
		{name: "remove category from dream", method: http.MethodDelete, url: "/dreams/1/categories/1", statusCode: http.StatusOK, response: emptyCategories},
		{name: "get categories", method: http.MethodGet, url: "/categories", statusCode: http.StatusOK, response: emptyCategories},
	}
	for _, test := range tests {
		test.evaluate(t, g)
	}
}

func TestMergeCategories(t *testing.T) {
	g := setupApiTest(t)
	date := time.Date(2020, 9, 13, 12, 30, 12, 0, time.UTC)

	var cat1 = entity.CategoryResponse{ID: 1, Name: "cat"}
	var cat2 = entity.CategoryResponse{ID: 2, Name: "new-cat"}
	var person1 = entity.CategoryResponse{ID: 3, Name: "person"}

	tests := []apiTest{
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: PostDreamRequest{Date: date}},
		{name: "add category to dream", method: http.MethodPut, url: "/dreams/1/categories?name=cat", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{cat1}}},
		{name: "add 2nd category to dream", method: http.MethodPut, url: "/dreams/1/categories?name=new-cat", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{cat1, cat2}}},
		{name: "post 2nd dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: PostDreamRequest{Date: date}},
		{name: "add person to 2nd dream", method: http.MethodPut, url: "/dreams/2/persons?name=person", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{cat1, cat2}, Persons: []entity.CategoryResponse{person1}}},
		{name: "merge categories in same dream", method: http.MethodPost, url: "/categories/merge", statusCode: http.StatusOK,
			body:     MergeCategoriesParams{SourceCategoryId: 1, TargetCategoryId: 2, NewName: "merged"},
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 2, Name: "merged"}}, Persons: []entity.CategoryResponse{person1}}},
		{name: "only 1 cat left in dream 1", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				Description: "", DreamMetaResponse: entity.DreamMetaResponse{
					ID: 1, Date: date, Finalized: false, Visible: true, Rating: nil,
					CategoriesResponse: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 2, Name: "merged"}}},
				}},
		},
		{name: "merge categories in different dreams", method: http.MethodPost, url: "/categories/merge", statusCode: http.StatusOK,
			body:     MergeCategoriesParams{SourceCategoryId: 2, TargetCategoryId: 3, NewName: "final merged"},
			response: entity.CategoriesResponse{Persons: []entity.CategoryResponse{{ID: 3, Name: "final merged"}}},
		},
		{name: "only 1 cat left in dream 1", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				Description: "", DreamMetaResponse: entity.DreamMetaResponse{
					ID: 1, Date: date, Finalized: false, Visible: true, Rating: nil,
					CategoriesResponse: entity.CategoriesResponse{Persons: []entity.CategoryResponse{{ID: 3, Name: "final merged"}}},
				}},
		},
	}
	for _, test := range tests {
		test.evaluate(t, g)
	}
}
