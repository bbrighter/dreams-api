package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAPITest(t *testing.T) (func(t *testing.T), *gin.Engine) {
	repo, deferedFunction := setupTest(t)
	var con Controller = InitController(repo)
	var router *gin.Engine = SetupRouter(con)

	return deferedFunction, router
}

func TestGetDreamsAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	handler := router.ServeHTTP
	assert.HTTPSuccess(t, handler, http.MethodGet, "/dreams", nil)
}

func TestCreateDreamAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	w := httptest.NewRecorder()
	params := `{"description": "desc", "date": "2021-02-18T21:54:42.123Z"}`
	paramsBytes := []byte(params)
	_bytes := bytes.NewReader(paramsBytes)
	req, _ := http.NewRequest(http.MethodPost, "/dreams", _bytes)

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Result().StatusCode)
}

func TestGetDreamAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	handler := router.ServeHTTP

	createTestDream(t)
	assert.HTTPSuccess(t, handler, http.MethodGet, "/dreams/1", nil)

	// Error handling
	// not found
	assert.HTTPStatusCode(t, handler, http.MethodGet, "/dreams/10", nil, 404)
	// bad param
	assert.HTTPStatusCode(t, handler, http.MethodGet, "/dreams/x", nil, 400)
}

func TestDeleteDreamAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)
	createTestDream(t)

	handler := router.ServeHTTP
	assert.HTTPSuccess(t, handler, http.MethodDelete, "/dreams/1", nil)

	// Error handling
	// not found
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/10", nil, 404)
	// bad param
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/x", nil, 400)
}
