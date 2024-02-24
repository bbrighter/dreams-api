package controller

import (
	"net/http"
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
	var statusCode int
	body = DreamRequestBody{Description: &desc, Date: time.Now()}

	statusCode = exectueTestRequest(http.MethodPost, "/dreams", body, nil, router)

	assert.Equal(t, 201, statusCode)

	// Missing params
	body = DreamRequestBody{Description: &desc}
	statusCode = exectueTestRequest(http.MethodPost, "/dreams", body, nil, router)
	assert.Equal(t, 400, statusCode)

	// Bad body
	badBody := `{"body":"bad"}`
	statusCode = exectueTestRequest(http.MethodPost, "/dreams", badBody, nil, router)
	assert.Equal(t, 400, statusCode)
}

func TestUpdateDreamAPI(t *testing.T) {
	teardown, router := setupAPITest(t)
	CreateOnlyTestDream(t)
	defer teardown(t)

	var body DreamRequestBody
	var desc string = "desc"
	var statusCode int
	body = DreamRequestBody{Description: &desc, Date: time.Now()}

	statusCode = exectueTestRequest(http.MethodPatch, "/dreams/1", body, nil, router)
	assert.Equal(t, 200, statusCode)

	// ID not uint
	statusCode = exectueTestRequest(http.MethodPatch, "/dreams/abc", body, nil, router)
	assert.Equal(t, 400, statusCode)

	// Not existing
	statusCode = exectueTestRequest(http.MethodPatch, "/dreams/100", body, nil, router)
	assert.Equal(t, 404, statusCode)

	// Bad body
	badBody := `{"hello": 1}`
	statusCode = exectueTestRequest(http.MethodPatch, "/dreams/1", badBody, nil, router)
	assert.Equal(t, 400, statusCode)
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

func TestGetPrivateDreams(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	status := exectueTestRequest(http.MethodGet, "/private/dreams", nil, nil, router)
	assert.Equal(t, 401, status)

	header := make(map[string]string)
	header["Authorization"] = "Basic dXNlcjowODAzODg="
	statusAuth := exectueTestRequest(http.MethodGet, "/private/dreams", nil, &header, router)
	assert.Equal(t, 200, statusAuth)
}

func TestGetPrivateDream(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	header := make(map[string]string)
	header["Authorization"] = "Basic dXNlcjowODAzODg="

	var statusCode int
	// Success
	statusCode = exectueTestRequest(http.MethodGet, "/private/dreams", nil, &header, router)
	assert.Equal(t, 200, statusCode)

	// Unauthorized
	statusCode = exectueTestRequest(http.MethodGet, "/private/dreams", nil, nil, router)
	assert.Equal(t, 401, statusCode)
}

func TestToggleVisiblity(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	header := make(map[string]string)
	header["Authorization"] = "Basic dXNlcjowODAzODg="

	var statusCode int

	// Valid
	createTestDream(t, 0, 0)
	statusCode = exectueTestRequest(http.MethodPatch, "/private/dreams/1", nil, &header, router)
	assert.Equal(t, 200, statusCode)

	// Not exists
	statusCode = exectueTestRequest(http.MethodPatch, "/private/dreams/100", nil, &header, router)
	assert.Equal(t, 404, statusCode)

	// Bad param
	statusCode = exectueTestRequest(http.MethodPatch, "/private/dreams/xy", nil, &header, router)
	assert.Equal(t, 400, statusCode)
}
