package controller

import (
	"net/http"
	"time"
)

func (s *ApiTestSuite) TestCategories() {
	date := time.Now()

	tests := []apiTest{
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated, body: PostDreamRequest{Date: date}},
		{name: "add category to dream", method: http.MethodPost, url: "/dreams/1/categories", statusCode: http.StatusCreated,
			body: PostCategoryRequestBody{Name: "cat", Type: "category"},
		},
		{name: "change type of category", method: http.MethodPatch, url: "/categories/1/type?type=person", statusCode: http.StatusOK},
		{name: "get categories", method: http.MethodGet, url: "/categories", statusCode: http.StatusOK,
			response: CategoryListResponse{Categories: []CategoryResponse{{ID: 1, Name: "cat", Type: "person"}}}},
		{name: "change name of category", method: http.MethodPatch, url: "/categories/1/name?name=newCat", statusCode: http.StatusOK},
		{name: "remove category from dream", method: http.MethodDelete, url: "/dreams/1/categories/1", statusCode: http.StatusOK},
		{name: "get categories", method: http.MethodGet, url: "/categories", statusCode: http.StatusOK,
			response: CategoryListResponse{Categories: []CategoryResponse{{ID: 1, Name: "newCat", Type: "person"}}}},
	}
	for _, test := range tests {
		s.evaluate(test)
	}
}

func (s *ApiTestSuite) TestMergeCategories() {
	date := time.Date(2020, 9, 13, 12, 30, 12, 0, time.UTC)

	tests := []apiTest{
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: PostDreamRequest{Date: date}},
		{name: "add category to dream", method: http.MethodPost, url: "/dreams/1/categories", statusCode: http.StatusCreated,
			body: PostCategoryRequestBody{Name: "cat", Type: "category"},
		},
		{name: "add 2nd category to dream", method: http.MethodPost, url: "/dreams/1/categories", statusCode: http.StatusCreated,
			body: PostCategoryRequestBody{Name: "new-cat", Type: "category"}},
		{name: "post 2nd dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: PostDreamRequest{Date: date}},
		{name: "add person to 2nd dream", method: http.MethodPost, url: "/dreams/2/categories", statusCode: http.StatusCreated,
			body: PostCategoryRequestBody{Name: "person", Type: "person"},
		},
		{name: "merge categories in same dream", method: http.MethodPost, url: "/categories/merge", statusCode: http.StatusOK,
			body: MergeCategoriesParams{SourceCategoryId: 1, TargetCategoryId: 2, NewName: "merged"},
			response: CategoryListResponse{Categories: []CategoryResponse{
				{ID: 2, Count: 1, Name: "merged", Type: "category"},
				{ID: 3, Count: 1, Name: "person", Type: "person"},
			}},
		},
		{name: "only 1 cat left in dream 1", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: DreamResponse{
				Categories:  []CategoryResponse{{ID: 2, Name: "merged", Type: "category"}},
				Description: "", ID: 1, Date: date, Finalized: false, Rating: nil},
		},
		{name: "merge categories in different dreams", method: http.MethodPost, url: "/categories/merge", statusCode: http.StatusOK,
			body:     MergeCategoriesParams{SourceCategoryId: 2, TargetCategoryId: 3, NewName: "final merged"},
			response: CategoryListResponse{Categories: []CategoryResponse{{ID: 3, Name: "final merged", Count: 2, Type: "person"}}},
		},
		{name: "only 1 cat left in dream 1", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: DreamResponse{
				Categories:  []CategoryResponse{{ID: 3, Name: "final merged", Type: "person"}},
				Description: "", ID: 1, Date: date, Finalized: false, Rating: nil},
		},
	}
	for _, test := range tests {
		s.evaluate(test)
	}
}
