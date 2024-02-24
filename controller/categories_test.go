package controller

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCategories(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	createTestDream(t, 1, 0)

	handler := router.ServeHTTP
	assert.HTTPSuccess(t, handler, http.MethodGet, "/categories", nil)
}

func TestAddCategory(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	id := CreateOnlyTestDream(t)

	handler := router.ServeHTTP
	validUrl := "/dreams/" + id + "/categories"
	values := url.Values{}
	values.Set("name", "Title")
	assert.HTTPSuccess(t, handler, http.MethodPut, validUrl, values)

	assert.HTTPStatusCode(t, handler, http.MethodPut, "/dreams/40/categories", values, 404)
	assert.HTTPStatusCode(t, handler, http.MethodPut, "/dreams/a/categories", values, 400)
	assert.HTTPStatusCode(t, handler, http.MethodPut, validUrl, nil, 400)
}

func TestRemoveCategoryOK(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)
	handler := router.ServeHTTP

	var dreamId, cateogryId string

	dreamId, categoryIds, _ := createTestDream(t, 1, 0)
	cateogryId = categoryIds[0]

	assert.HTTPSuccess(t, handler, http.MethodDelete, "/dreams/"+dreamId+"/categories/"+cateogryId, nil)
}
func TestRemoveCategoryErrors(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)
	handler := router.ServeHTTP

	var dreamId, categoryId string

	dreamId, categoryIds, _ := createTestDream(t, 1, 0)
	categoryId = categoryIds[0]

	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/100/categories/"+categoryId, nil, 404)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/"+dreamId+"/categories/100", nil, 404)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/a/categories/1", nil, 400)
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/1/categories/a", nil, 400)
}
