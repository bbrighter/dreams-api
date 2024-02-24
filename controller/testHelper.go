package controller

import (
	"bytes"
	"encoding/json"
	"io"
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
	dream := store.CreateTestDream(numberOfCategories, numberOfPersons, true, t)
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

func exectueTestRequest(
	method string,
	url string,
	body interface{},
	header *map[string]string,
	router *gin.Engine,
) int {
	var bodyBytes io.Reader = nil
	if body != nil {
		requestBody, _ := json.Marshal(body)
		bodyBytes = bytes.NewBuffer(requestBody)
	}
	request, _ := http.NewRequest(method, url, bodyBytes)

	if header != nil {
		for key, value := range *header {
			request.Header.Add(key, value)
		}
	}

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	return writer.Result().StatusCode
}
