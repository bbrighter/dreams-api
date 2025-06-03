package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

func TestPersons(t *testing.T) {
	g := setupApiTest(t)
	date := time.Now()

	tests := []apiTest{
		{name: "post dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
			body: DreamRequestBody{Date: date}},
		{name: "add person to dream", method: http.MethodPut, url: "/dreams/1/persons?name=person", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{{ID: 1, Name: "person"}}}},
		{name: "get persons", method: http.MethodGet, url: "/persons", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{{ID: 1, Name: "person"}}}},
		{name: "remove person from dream", method: http.MethodDelete, url: "/dreams/1/persons/1", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{}}},
		{name: "get persons", method: http.MethodGet, url: "/persons", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{}}},
	}
	for _, test := range tests {
		test.evaluate(t, g)
	}
}
