package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestRoutePrivateDreams() (*privateDreamsRoute, *gin.Context, *httptest.ResponseRecorder) {
	var tuc = testUseCaseDreams{}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	return &privateDreamsRoute{d: tuc}, c, rec
}

func TestGetAllPrivate(t *testing.T) {
	r, g, rec := newTestRoutePrivateDreams()

	r.GetAll(g)

	assert.Equal(t, 200, rec.Code)
}

func TestGetPrivate(t *testing.T) {
	r, g, rec := newTestRoutePrivateDreams()
	g.AddParam("id", "1")

	r.Get(g)

	assert.Equal(t, 200, rec.Code)
}

func TestGetPrivate404(t *testing.T) {
	r, g, rec := newTestRoutePrivateDreams()
	g.AddParam("id", "2")

	r.Get(g)

	assert.Equal(t, 404, rec.Code)
}

func TestToggleVisiblity(t *testing.T) {
	r, g, rec := newTestRoutePrivateDreams()
	g.AddParam("id", "1")

	r.ToggleVisibility(g)

	assert.Equal(t, rec.Code, 200)
}
