package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

func TestDreams(t *testing.T) {
	h := setupApiTest(t)

	desc := "description"
	date := time.Date(2022, 11, 13, 4, 12, 8, 0, time.UTC)

	tests := []apiTest{
		{name: "get dreams, empty", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{}}},
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: DreamRequestBody{Date: date}},
		{name: "get dreams, one exits", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{{ID: 1, Date: date, Finalized: false, Visible: true}}}},
		{name: "get dreams with persons and categories", method: http.MethodGet, url: "/dreams?includes=persons,categories", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{{
				ID: 1, Date: date, Finalized: false, Visible: true,
			}}}},
		{name: "get single dream", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				Description:       "",
				DreamMetaResponse: entity.DreamMetaResponse{ID: 1, Date: date, Finalized: false, Visible: true},
			},
		},
		{name: "patch dream", method: http.MethodPatch, url: "/dreams/1", statusCode: http.StatusOK,
			body: DreamRequestBody{Description: &desc, Date: date}},
		{name: "add person to dream", method: http.MethodPut, url: "/dreams/1/persons?name=person", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{{ID: 1, Name: "person"}}}},
		{name: "add category to dream", method: http.MethodPut, url: "/dreams/1/categories?name=cat", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 1, Name: "cat"}}}},
		{name: "get dreams with persons and categories", method: http.MethodGet, url: "/dreams?includes=persons,categories", statusCode: http.StatusOK,
			response: entity.DreamsResponse{
				Dreams: []entity.DreamMetaResponse{{
					ID: 1, Date: date, Finalized: false, Visible: true,
					Categories: []entity.CategoryResponse{{ID: 1, Name: "cat"}},
					Persons:    []entity.PersonResponse{{ID: 1, Name: "person"}}},
				}}},
		{name: "remove person to dream", method: http.MethodDelete, url: "/dreams/1/persons/1", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{}}},
		{name: "remove category to dream", method: http.MethodDelete, url: "/dreams/1/categories/1", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{}}},
		{name: "Finalize the dream", method: http.MethodPatch, url: "/dreams/1/finalize", statusCode: http.StatusOK},
		{name: "delete dream", method: http.MethodDelete, url: "/dreams/1", statusCode: http.StatusOK},
	}
	for _, test := range tests {
		test.evaluate(t, h)
	}
}

func TestPrivateDreams(t *testing.T) {
	h := setupApiTest(t)

	date := time.Date(2022, 11, 13, 4, 12, 8, 0, time.UTC)

	tests := []apiTest{
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated, body: DreamRequestBody{Date: date}},
		{name: "toggle visibility of dream", method: http.MethodPatch, url: "/dreams/private/1", statusCode: http.StatusOK},
		{name: "private dreams not listed", method: http.MethodGet, url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{}}},
		{name: "private dreams listed", method: http.MethodGet, url: "/dreams/private", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{{ID: 1, Date: date, Finalized: false, Visible: false}}}},
		{name: "private single dream not found", method: http.MethodGet, url: "/dreams/1", statusCode: http.StatusNotFound},
		{name: "get private single dream", method: http.MethodGet, url: "/dreams/private/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				Description:       "",
				DreamMetaResponse: entity.DreamMetaResponse{ID: 1, Date: date, Finalized: false, Visible: false}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.evaluate(t, h)
		})
	}
}
