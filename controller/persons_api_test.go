package controller

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPersons(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	createTestDream(t, 0, 10)

	handler := router.ServeHTTP

	assert.HTTPSuccess(t, handler, http.MethodGet, "/persons", nil)
}

func TestPutPersonToDream(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	dreamId := CreateOnlyTestDream(t)

	handler := router.ServeHTTP
	validUrl := "/dreams/" + dreamId + "/persons"
	values := url.Values{}
	values.Set("name", "Name")

	assert.HTTPSuccess(t, handler, http.MethodPut, validUrl, values)
	assert.HTTPStatusCode(t, handler, http.MethodPut, "/dreams/40/persons", values, 404)
	assert.HTTPStatusCode(t, handler, http.MethodPut, "/dreams/a/persons", values, 400)
	assert.HTTPStatusCode(t, handler, http.MethodPut, validUrl, nil, 400)
}

func TestRemovePersonFromDream(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	handler := router.ServeHTTP
	dreamId, _, personIds := createTestDream(t, 0, 2)
	validUrl := "/dreams/" + dreamId + "/persons/" + personIds[0]

	assert.HTTPSuccess(t, handler, http.MethodDelete, validUrl, nil)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/100/persons/1", nil, 404)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/1/persons/100", nil, 404)
}
