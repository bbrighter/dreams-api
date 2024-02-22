package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDreamsAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	handler := router.ServeHTTP
	assert.HTTPSuccess(t, handler, http.MethodGet, "/dreams", nil)
}

func TestCreateDreamAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	var body DreamRequestBody
	var desc string = "desc"
	var w *httptest.ResponseRecorder
	body = DreamRequestBody{Description: &desc, Date: time.Now()}

	w = makeRequest(http.MethodPost, "/dreams", body, router)

	assert.Equal(t, 201, w.Result().StatusCode)

	// Missing params
	body = DreamRequestBody{Description: &desc}
	w = makeRequest(http.MethodPost, "/dreams", body, router)
	assert.Equal(t, 400, w.Result().StatusCode)

	body = DreamRequestBody{Date: time.Now()}
	w = makeRequest(http.MethodPost, "/dreams", body, router)
	assert.Equal(t, 201, w.Result().StatusCode)
}

func TestUpdateDreamAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	CreateOnlyTestDream(t)
	defer teardown(t)

	var body DreamRequestBody
	var desc string = "desc"
	var w *httptest.ResponseRecorder
	body = DreamRequestBody{Description: &desc, Date: time.Now()}

	w = makeRequest(http.MethodPatch, "/dreams/1", body, router)
	assert.Equal(t, 200, w.Result().StatusCode)

	visible := false
	body = DreamRequestBody{Description: &desc, Date: time.Now(), Visible: &visible}
	w = makeRequest(http.MethodPatch, "/dreams/1", body, router)
	assert.Equal(t, 200, w.Result().StatusCode)

	// ID not uint
	w = makeRequest(http.MethodPatch, "/dreams/abc", body, router)
	assert.Equal(t, 400, w.Result().StatusCode)
}

func TestGetDreamAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	handler := router.ServeHTTP

	CreateOnlyTestDream(t)
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
	CreateOnlyTestDream(t)

	handler := router.ServeHTTP
	assert.HTTPSuccess(t, handler, http.MethodDelete, "/dreams/1", nil)

	// Error handling
	// not found
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/10", nil, 404)
	// bad param
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/x", nil, 400)
}
