package controller

import (
	"net/http"
	"time"
)

func (s *ApiTestSuite) TestDreams() {
	var catResp = CategoryListResponse{Categories: []CategoryResponse{
		{ID: 1, Name: "person", Type: "person"},
		{ID: 2, Name: "cat", Type: "category"}},
	}

	desc := "description"
	date := time.Date(2022, 11, 13, 4, 12, 8, 0, time.UTC)
	rating := 3

	tests := []apiTest{
		{name: "get dreams, empty", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: DreamListResponse{Dreams: []DreamResponse{}}},
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: PostDreamRequest{Date: date}},
		{name: "get dreams, one exits", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: DreamListResponse{
				Dreams: []DreamResponse{{ID: 1, Date: date, Finalized: false, Categories: []CategoryResponse{}}}}},
		{name: "get single dream", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: DreamResponse{
				Description: "", ID: 1, Date: date, Finalized: false, Categories: []CategoryResponse{},
			},
		},
		{name: "patch dream", method: http.MethodPatch, url: "/dreams/1", statusCode: http.StatusOK, body: UpdateDreamRequest{Description: &desc, Date: &date}},
		{name: "patch dream only date", method: http.MethodPatch, url: "/dreams/1", statusCode: http.StatusOK, body: UpdateDreamRequest{Date: &date}},
		{name: "add person to dream", method: http.MethodPost, url: "/dreams/1/categories", statusCode: http.StatusCreated, body: PostCategoryRequestBody{Name: "person", Type: "person"}},
		{name: "add category to dream", method: http.MethodPost, url: "/dreams/1/categories", statusCode: http.StatusCreated, body: PostCategoryRequestBody{Name: "cat", Type: "category"}},
		{name: "get single dream with all changes made", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: DreamResponse{
				Description: desc,
				ID:          1,
				Date:        date,
				Finalized:   false,
				Categories:  catResp.Categories,
			}},
		{name: "list dreams, one exists", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: DreamListResponse{
				Dreams: []DreamResponse{{
					ID:          1,
					Description: desc,
					Date:        date,
					Finalized:   false,
					Categories:  catResp.Categories,
				}}}},
		{name: "remove person from dream", method: http.MethodDelete, url: "/dreams/1/categories/1", statusCode: http.StatusOK},
		{name: "remove category from dream", method: http.MethodDelete, url: "/dreams/1/categories/2", statusCode: http.StatusOK},
		{name: "rate the dream", method: http.MethodPatch, url: "/dreams/1", statusCode: http.StatusOK,
			body: UpdateDreamRequest{Rating: &rating}},
		{name: "get rated dream", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: DreamResponse{
				Description: desc, ID: 1, Date: date, Finalized: false, Rating: &rating, Categories: []CategoryResponse{},
			}},
		{name: "finalize dream", method: http.MethodPatch, url: "/dreams/1/finalize", statusCode: http.StatusOK},
		{name: "dream is finalized", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: DreamResponse{Description: desc, ID: 1, Date: date, Finalized: true, Rating: &rating, Categories: []CategoryResponse{}},
		},
		{name: "get dreams with rating and finalization", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: DreamListResponse{
				Dreams: []DreamResponse{{
					ID: 1, Date: date, Finalized: true, Rating: &rating, Description: "description",
					Categories: []CategoryResponse{},
				},
				}}},
		{name: "delete dream", method: http.MethodDelete, url: "/dreams/1", statusCode: http.StatusOK},
	}
	for _, test := range tests {
		s.evaluate(test)
	}
}

func (s *ApiTestSuite) TestBadParams() {

	date := time.Date(2022, 11, 13, 4, 12, 8, 0, time.UTC)
	var zeroTime time.Time
	desc := "desc"

	tests := []apiTest{
		{name: "post dream no params", method: http.MethodPost, url: "/dreams", statusCode: http.StatusBadRequest, body: PostDreamRequest{}},
		{name: "post dream only desc", method: http.MethodPost, url: "/dreams", statusCode: http.StatusBadRequest, body: PostDreamRequest{Description: &desc}},
		{name: "post dream with zero date", method: http.MethodPost, url: "/dreams", statusCode: http.StatusBadRequest, body: PostDreamRequest{Date: zeroTime}},
		{name: "get dream not found", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusNotFound},
		{name: "get dream bad param", method: http.MethodGet, url: "/dreams/x", statusCode: http.StatusBadRequest},
		{name: "post dream successfully", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated, body: PostDreamRequest{Date: date}},
		{name: "finalize dream without rating", method: http.MethodPatch, url: "/dreams/1/finalize", statusCode: http.StatusBadRequest},
		{name: "patch dream bad param", method: http.MethodPatch, url: "/dreams/1", statusCode: http.StatusBadRequest, body: UpdateDreamRequest{}},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			s.evaluate(test)
		})
	}
}
