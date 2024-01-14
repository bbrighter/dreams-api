package controller

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	createTestTagAndDream(t)

	handler := router.ServeHTTP
	assert.HTTPSuccess(t, handler, http.MethodGet, "/tags", nil)
}

func TestAddTag(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	id := createTestDream(t)

	handler := router.ServeHTTP
	validUrl := "/dreams/" + strconv.FormatUint(uint64(id), 10) + "/tags"
	values := url.Values{}
	values.Set("title", "Title")
	assert.HTTPSuccess(t, handler, http.MethodPut, validUrl, values)

	assert.HTTPStatusCode(t, handler, http.MethodPut, "/dreams/40/tags", values, 404)
	assert.HTTPStatusCode(t, handler, http.MethodPut, "/dreams/a/tags", values, 400)
	assert.HTTPStatusCode(t, handler, http.MethodPut, validUrl, nil, 400)
}

func TestRemoveTag(t *testing.T) {
	validUrl := func(dreamId uint, tagId uint) string {
		return "/dreams/" + strconv.FormatUint(uint64(dreamId), 10) +
			"/tags/" + strconv.FormatUint(uint64(tagId), 10)
	}

	teardown, router := setupAPITest(t)
	defer teardown(t)
	handler := router.ServeHTTP

	var dreamId uint
	var tagId uint

	dreamId, tagId = createTestTagAndDream(t)

	assert.HTTPSuccess(t, handler, http.MethodDelete, validUrl(dreamId, tagId), nil)

	cleanTestEntries(t)
	dreamId, tagId = createTestTagAndDream(t)

	assert.HTTPStatusCode(t, handler, http.MethodDelete, validUrl(100, tagId), nil, 404)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, validUrl(dreamId, 100), nil, 404)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/a/tags/1", nil, 400)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/1/tags/a", nil, 400)
}
