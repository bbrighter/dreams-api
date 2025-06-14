package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

func TestGetStatistics(t *testing.T) {
	g := setupApiTest(t)

	tests := []apiTest{
		{name: "get empty statistics", method: http.MethodGet, url: "/statistics?limit=40", statusCode: http.StatusOK,
			response: entity.CountsResponse{Categories: []entity.CountResponse{}, Persons: []entity.CountResponse{}}},
		{name: "setup dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: DreamRequestBody{Date: time.Now()}},
		{name: "setup person", method: http.MethodPut, url: "/dreams/1/persons?name=Pers", statusCode: http.StatusOK},
		{name: "setup category", method: http.MethodPut, url: "/dreams/1/categories?name=Cat", statusCode: http.StatusOK},
		{name: "get filled statistics", method: http.MethodGet, url: "/statistics?limit=40", statusCode: http.StatusOK,
			response: entity.CountsResponse{
				Categories: []entity.CountResponse{{ID: 2, Count: 1}},
				Persons:    []entity.CountResponse{{ID: 1, Count: 1}},
			}},
	}

	for _, test := range tests {
		test.evaluate(t, g)
	}
}
