package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func createTestDream(t *testing.T, numberOfCategories int, numberOfPersons int) (string, []string, []string) {
	dream := store.CreateTestDream(numberOfCategories, numberOfPersons, t)
	dreamId := strconv.FormatUint(uint64(dream.ID), 10)
	categoryIds := []string{}
	for _, c := range dream.Categories {
		categoryIds = append(categoryIds, strconv.FormatUint(uint64(c.ID), 10))
	}
	personIds := []string{}
	for _, p := range dream.Persons {
		personIds = append(personIds, strconv.FormatUint(uint64(p.ID), 10))
	}
	return dreamId, categoryIds, personIds
}

func CreateOnlyTestDream(t *testing.T) string {
	id, _, _ := createTestDream(t, 0, 0)
	return id
}

func makeRequest(method, url string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	return writer
}
