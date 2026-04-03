package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

func (s *ApiTestSuite) TestDreams() {

	desc := "description"
	date := time.Date(2022, 11, 13, 4, 12, 8, 0, time.UTC)
	rating := 3

	tests := []apiTest{
		{name: "get dreams, empty", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamResponse{}}},
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: PostDreamRequest{Date: date}},
		{name: "get dreams, one exits", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{
				Dreams: []entity.DreamResponse{{ID: 1, Date: date, Finalized: false, Visible: true, CategoriesResponse: entity.CategoriesResponse{}}}}},
		{name: "get dreams with persons and categories", method: http.MethodGet, url: "/dreams?includes=persons,categories", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamResponse{{
				ID: 1, Date: date, Finalized: false, Visible: true,
			}}}},
		{name: "get single dream", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				Description: "", ID: 1, Date: date, Finalized: false, Visible: true,
			},
		},
		{name: "patch dream", method: http.MethodPatch, url: "/dreams/1", statusCode: http.StatusOK,
			body: UpdateDreamRequest{Description: &desc, Date: &date}},
		{name: "patch dream only date", method: http.MethodPatch, url: "/dreams/1", statusCode: http.StatusOK,
			body: UpdateDreamRequest{Date: &date}},
		{name: "add person to dream", method: http.MethodPost, url: "/dreams/1/categories", statusCode: http.StatusCreated,
			body: PostCategoryRequestBody{Name: "person", Type: entity.TypePerson},
		},
		{name: "add category to dream", method: http.MethodPost, url: "/dreams/1/categories", statusCode: http.StatusCreated,
			body: PostCategoryRequestBody{Name: "cat", Type: entity.TypeCategory},
		},
		{name: "get single dream with all changes made", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				Description: desc, ID: 1, Date: date, Finalized: false, Visible: true,
				CategoriesResponse: entity.CategoriesResponse{
					Categories: []entity.CategoryResponse{
						{ID: 1, Name: "person", Type: entity.TypePerson},
						{ID: 2, Name: "cat", Type: entity.TypeCategory}},
				},
			}},
		{name: "remove person from dream", method: http.MethodDelete, url: "/dreams/1/categories/1", statusCode: http.StatusOK},
		{name: "remove category from dream", method: http.MethodDelete, url: "/dreams/1/categories/2", statusCode: http.StatusOK},
		{name: "rate the dream", method: http.MethodPatch, url: "/dreams/1", statusCode: http.StatusOK,
			body: UpdateDreamRequest{Rating: &rating}},
		{name: "get rated dream", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				Description: desc, ID: 1, Date: date, Finalized: false, Visible: true, Rating: &rating,
			}},
		{name: "Finalize the dream", method: http.MethodPatch, url: "/dreams/1/finalize", statusCode: http.StatusOK},
		{name: "Dream is finalized", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{Description: desc, ID: 1, Date: date, Finalized: true, Visible: true, Rating: &rating},
		},
		{name: "Get dreams with rating and finalization", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{
				Dreams: []entity.DreamResponse{{
					ID: 1, Date: date, Finalized: true, Visible: true, Rating: &rating, Description: "description",
					CategoriesResponse: entity.CategoriesResponse{},
				},
				}}},
		{name: "delete dream", method: http.MethodDelete, url: "/dreams/1", statusCode: http.StatusOK},
	}
	for _, test := range tests {
		s.evaluate(test)
	}
}

func (s *ApiTestSuite) TestPrivateDreams() {
	date := time.Date(2022, 11, 13, 4, 12, 8, 0, time.UTC)
	var token *string
	extractToken := func(resp *httptest.ResponseRecorder) {
		var loginResp entity.LoginResponse
		err := json.NewDecoder(resp.Body).Decode(&loginResp)
		s.Require().NoError(err)
		token = &loginResp.Token
	}
	tokenFn := func() *string { return token }

	tests := []apiTest{
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated, body: PostDreamRequest{Date: date}},
		{name: "login", method: http.MethodPost, url: "/login", statusCode: http.StatusOK, body: LoginRequest{Name: "Benni", Password: "0803"}, after: extractToken},
		{name: "toggle visibility of dream", method: http.MethodPatch, url: "/dreams/private/1", statusCode: http.StatusOK, token: tokenFn},
		{name: "private dreams not listed", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamResponse{}}},
		{name: "private dreams listed", method: http.MethodGet, url: "/dreams/private", statusCode: http.StatusOK, token: tokenFn,
			response: entity.DreamsResponse{Dreams: []entity.DreamResponse{{ID: 1, Date: date, Finalized: false, Visible: false}}}},
		{name: "private single dream not found", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusNotFound},
		{name: "get private single dream", method: http.MethodGet, url: "/dreams/private/1", statusCode: http.StatusOK, token: tokenFn,
			response: entity.DreamResponse{
				Description: "", ID: 1, Date: date, Finalized: false, Visible: false},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			s.evaluate(test)
		})
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
