package controller

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	createTestDream(t, 1, 0)

	handler := router.ServeHTTP
	assert.HTTPSuccess(t, handler, http.MethodGet, "/tags", nil)
}

func TestAddTag(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	id := CreateOnlyTestDream(t)

	handler := router.ServeHTTP
	validUrl := "/dreams/" + id + "/tags"
	values := url.Values{}
	values.Set("title", "Title")
	assert.HTTPSuccess(t, handler, http.MethodPut, validUrl, values)

	assert.HTTPStatusCode(t, handler, http.MethodPut, "/dreams/40/tags", values, 404)
	assert.HTTPStatusCode(t, handler, http.MethodPut, "/dreams/a/tags", values, 400)
	assert.HTTPStatusCode(t, handler, http.MethodPut, validUrl, nil, 400)
}

func TestRemoveTagOK(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)
	handler := router.ServeHTTP

	var dreamId, tagId string

	dreamId, tagIds, _ := createTestDream(t, 1, 0)
	tagId = tagIds[0]

	assert.HTTPSuccess(t, handler, http.MethodDelete, "/dreams/"+dreamId+"/tags/"+tagId, nil)
}
func TestRemoveTagErrors(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)
	handler := router.ServeHTTP

	var dreamId, tagId string

	dreamId, tagIds, _ := createTestDream(t, 1, 0)
	tagId = tagIds[0]

	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/100/tags/"+tagId, nil, 404)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/"+dreamId+"/tags/100", nil, 404)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/a/tags/1", nil, 400)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/1/tags/a", nil, 400)
}
