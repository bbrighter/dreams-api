package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bbrighter/dreams-api/store"
	"github.com/gin-gonic/gin"
)

func setupAPITest(t *testing.T) (func(t *testing.T), *gin.Engine) {
	repo, deferedFunction := store.SetupTest(t)
	var con Controller = InitController(repo)
	var router *gin.Engine = SetupRouter(con)

	return deferedFunction, router
}

func createTestDream(t *testing.T) uint {
	dream := store.CreateTestDream(t)
	return dream.ID
}

func createTestTagAndDream(t *testing.T) (uint, uint) {
	dream := store.CreateTestTagAndDream(t)
	return dream.ID, dream.Tags[0].ID
}

func cleanTestEntries(t *testing.T) {
	store.CleanTestEntries(t)
}

func makeRequest(method, url string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	return writer
}
