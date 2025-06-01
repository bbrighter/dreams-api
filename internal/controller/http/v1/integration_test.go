package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/bbrighter/dreams-api/internal/usecase/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func setupApiTest(t *testing.T) *gin.Engine {
	db := setupTestDB()
	logger := zap.L()
	dreamsRepo := repository.NewDreamsRepo(db)
	dreamsUseCase := usecase.NewDreamUseCase(dreamsRepo)
	privateDreamsUseCase := usecase.NewPrivateDreamUseCase(dreamsRepo)
	personsRepo := repository.NewPersonsRepo(db)
	personsUseCase := usecase.NewPersonsUseCase(personsRepo)
	categoriesRepo := repository.NewCategoriesRepo(db)
	categoriesUseCase := usecase.NewCategoriesUseCase(categoriesRepo)
	statisticsRepo := repository.NewStatisticsRepo(db)
	statisticsUseCase := usecase.NewStatisticsUseCase(statisticsRepo)
	repository.Migration(db, logger)

	gin.SetMode(gin.TestMode)
	handler := gin.New()
	NewRouter(
		handler,
		dreamsUseCase,
		privateDreamsUseCase,
		personsUseCase,
		categoriesUseCase,
		statisticsUseCase,
		categoriesUseCase,
		personsUseCase,
	)
	return handler
}

func TestDreams(t *testing.T) {
	h := setupApiTest(t)

	desc := "description"
	date := time.Date(2022, 11, 13, 4, 12, 8, 0, time.UTC)

	tests := []struct {
		name       string
		method     string
		url        string
		body       any
		statusCode int
		response   any
	}{
		{name: "get dreams, empty", method: "GET", url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{}}},
		{name: "post dream", method: "POST", url: "/dreams", statusCode: http.StatusCreated,
			body: DreamRequestBody{Date: date}},
		{name: "get dreams, one exits", method: "GET", url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{{ID: 1, Date: date, Finalized: false, Visible: true}}}},
		{name: "get single dream", method: "GET", url: "/dreams/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				DreamMetaResponse: entity.DreamMetaResponse{ID: 1, Date: date, Finalized: false, Visible: true},
				Description:       "",
				Categories:        entity.CategoriesResponse{Categories: []entity.CategoryResponse{}},
				Persons:           entity.PersonsResponse{Persons: []entity.PersonResponse{}}},
		},
		{name: "patch dream", method: "PATCH", url: "/dreams/1", statusCode: http.StatusOK,
			body: DreamRequestBody{Description: &desc, Date: date}},
		{name: "add person to dream", method: "PUT", url: "/dreams/1/persons?name=person", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{{ID: 1, Name: "person"}}}},
		{name: "remove person to dream", method: "DELETE", url: "/dreams/1/persons/1", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{}}},
		{name: "add category to dream", method: "PUT", url: "/dreams/1/categories?name=cat", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 1, Name: "cat"}}}},
		{name: "remove category to dream", method: "DELETE", url: "/dreams/1/categories/1", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{}}},
		{name: "Finalize the dream", method: "PATCH", url: "/dreams/1/finalize", statusCode: http.StatusOK},
		{name: "delete dream", method: "DELETE", url: "/dreams/1", statusCode: http.StatusOK},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body io.Reader
			if test.body != nil {
				marBody, _ := json.Marshal(test.body)
				body = bytes.NewBuffer(marBody)
			}
			req, _ := http.NewRequest(test.method, test.url, body)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)
			assert.Equal(t, test.statusCode, resp.Code)

			if test.response != nil {
				expected := test.response
				actual := reflect.New(reflect.TypeOf(expected)).Interface()
				err := json.Unmarshal(resp.Body.Bytes(), &actual)
				assert.NoError(t, err)
				assert.Equal(t, expected, reflect.Indirect(reflect.ValueOf(actual)).Interface())
			}
		})
	}
}

func TestCategoriesAndPersons(t *testing.T) {
	g := setupApiTest(t)
	date := time.Now()

	tests := []struct {
		name       string
		method     string
		url        string
		body       any
		statusCode int
		response   any
	}{
		{name: "post dream", method: "POST", url: "/dreams", statusCode: http.StatusCreated,
			body: DreamRequestBody{Date: date}},
		{name: "add person to dream", method: "PUT", url: "/dreams/1/persons?name=person", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{{ID: 1, Name: "person"}}}},
		{name: "add category to dream", method: "PUT", url: "/dreams/1/categories?name=cat", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 1, Name: "cat"}}}},
		{name: "get categories", method: "GET", url: "/categories", statusCode: http.StatusOK,
			response: entity.CategoriesResponse{Categories: []entity.CategoryResponse{{ID: 1, Name: "cat"}}}},
		{name: "get persons", method: "GET", url: "/persons", statusCode: http.StatusOK,
			response: entity.PersonsResponse{Persons: []entity.PersonResponse{{ID: 1, Name: "person"}}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body io.Reader
			if test.body != nil {
				marBody, _ := json.Marshal(test.body)
				body = bytes.NewBuffer(marBody)
			}
			req, _ := http.NewRequest(test.method, test.url, body)
			resp := httptest.NewRecorder()
			g.ServeHTTP(resp, req)
			assert.Equal(t, test.statusCode, resp.Code)

			if test.response != nil {
				expected := test.response
				actual := reflect.New(reflect.TypeOf(expected)).Interface()
				err := json.Unmarshal(resp.Body.Bytes(), &actual)
				assert.NoError(t, err)
				assert.Equal(t, expected, reflect.Indirect(reflect.ValueOf(actual)).Interface())
			}
		})
	}
}

func TestPrivateDreams(t *testing.T) {
	h := setupApiTest(t)

	date := time.Date(2022, 11, 13, 4, 12, 8, 0, time.UTC)

	tests := []struct {
		name       string
		method     string
		url        string
		body       any
		statusCode int
		response   any
	}{
		{name: "post dream", method: "POST", url: "/dreams", statusCode: http.StatusCreated, body: DreamRequestBody{Date: date}},
		{name: "toggle visibility of dream", method: "PATCH", url: "/dreams/private/1", statusCode: http.StatusOK},
		{name: "private dreams not listed", method: "GET", url: "/dreams", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{}}},
		{name: "private dreams listed", method: "GET", url: "/dreams/private", statusCode: http.StatusOK,
			response: entity.DreamsResponse{Dreams: []entity.DreamMetaResponse{{ID: 1, Date: date, Finalized: false, Visible: false}}}},
		{name: "private single dream not found", method: "GET", url: "/dreams/1", statusCode: http.StatusNotFound},
		{name: "get private single dream", method: "GET", url: "/dreams/private/1", statusCode: http.StatusOK,
			response: entity.DreamResponse{
				DreamMetaResponse: entity.DreamMetaResponse{ID: 1, Date: date, Finalized: false, Visible: false},
				Description:       "",
				Categories:        entity.CategoriesResponse{Categories: []entity.CategoryResponse{}},
				Persons:           entity.PersonsResponse{Persons: []entity.PersonResponse{}}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body io.Reader
			if test.body != nil {
				marBody, _ := json.Marshal(test.body)
				body = bytes.NewBuffer(marBody)
			}
			req, _ := http.NewRequest(test.method, test.url, body)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)
			assert.Equal(t, test.statusCode, resp.Code)

			if test.response != nil {
				expected := test.response
				actual := reflect.New(reflect.TypeOf(expected)).Interface()
				err := json.Unmarshal(resp.Body.Bytes(), &actual)
				assert.NoError(t, err)
				assert.Equal(t, expected, reflect.Indirect(reflect.ValueOf(actual)).Interface())
			}
		})
	}
}
