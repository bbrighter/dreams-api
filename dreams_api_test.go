package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAPITest(t *testing.T) (func(t *testing.T), *gin.Engine) {
	repo, deferedFunction := setupTest(t)
	var con Controller = InitController(repo)
	var router *gin.Engine = SetupRouter(con)

	return deferedFunction, router
}

func makeRequest(method, url string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	return writer
}

func TestDreamRequestBodyToDream(t *testing.T) {
	t.Parallel()
	var body DreamRequestBody
	var dream Dream
	var desc string = "desc"

	body = DreamRequestBody{
		Date:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: &desc,
	}
	dream = body.dreamRequestBodyToDream()

	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), dream.Date)
	assert.Equal(t, "desc", dream.Description)

	body = DreamRequestBody{
		Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	dream = body.dreamRequestBodyToDream()
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), dream.Date)
	assert.Equal(t, "", dream.Description)

}

func TestDreamToDreamResponse(t *testing.T) {
	t.Parallel()
	var dream Dream = Dream{
		ID:          1,
		Date:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "desc",
	}

	//
	var resp DreamResponse = dreamToDreamResponse(dream)
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, "desc", resp.Description)
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), resp.Date)

}

func TestDreamsToDreamsResponse(t *testing.T) {
	t.Parallel()

	var dreams []Dream

	// Empty input gives {dreams: []}
	var emptyResp DreamsResponse = dreamsToDreamsResponse(dreams)
	assert.NotNil(t, emptyResp.Dreams)
	assert.Len(t, emptyResp.Dreams, 0)

	// Non-empty input gives {dreams: [...]}
	dreams = []Dream{
		{ID: 1, Date: time.Now(), Description: "desc"},
	}
	var resp DreamsResponse = dreamsToDreamsResponse(dreams)
	assert.Len(t, resp.Dreams, 1)
}

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
	createTestDream(t)
	defer teardown(t)

	var body DreamRequestBody
	var desc string = "desc"
	var w *httptest.ResponseRecorder
	body = DreamRequestBody{Description: &desc, Date: time.Now()}

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

	createTestDream(t)
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
	createTestDream(t)

	handler := router.ServeHTTP
	assert.HTTPSuccess(t, handler, http.MethodDelete, "/dreams/1", nil)

	// Error handling
	// not found
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/10", nil, 404)
	// bad param
	assert.HTTPStatusCode(t, handler, http.MethodDelete, "/dreams/x", nil, 400)
}
