package controller

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCountCategories(t *testing.T) {
	teardown, router := setupAPITest(t)
	defer teardown(t)

	createTestDream(t, 10, 0)
	handler := router.ServeHTTP

	assert.HTTPSuccess(t, handler, http.MethodGet, "/statistics", nil)

}
