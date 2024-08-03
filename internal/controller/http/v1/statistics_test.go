package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testUseCaseStatistics struct{}

func (t testUseCaseStatistics) GetStatistics(showAll bool, maxNumber int) (entity.Counts, entity.Counts) {
	var catCounts = entity.Counts{entity.Count{ID: 1, Count: 10}}
	var perCounts = entity.Counts{entity.Count{ID: 1, Count: 1}, entity.Count{ID: 2, Count: 5}}
	return catCounts, perCounts
}

func newTestRouteStatistics() (*statisticsRoute, *gin.Context, *httptest.ResponseRecorder) {
	var tuc = testUseCaseStatistics{}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	return &statisticsRoute{c: tuc}, c, rec
}

func TestGetStatistics(t *testing.T) {
	r, g, rec := newTestRouteStatistics()

	r.GetStatistics(g)
	assert.Equal(t, rec.Code, 200)
	expectedResp := `{"categories":[{"id":1,"count":10}],"persons":[{"id":1,"count":1},{"id":2,"count":5}]}`
	assert.Equal(t, rec.Body.String(), expectedResp)
}
