package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/bbrighter/dreams-api/internal/usecase/repository"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	return db
}

func setupApiTest(t *testing.T) *gin.Engine {
	logger := zap.NewExample()
	db := setupTestDB(t)
	dreamsRepo := repository.NewDreamsRepo(db)
	dreamsUseCase := usecase.NewDreamUseCase(dreamsRepo, logger)
	privateDreamsUseCase := usecase.NewPrivateDreamUseCase(dreamsRepo)
	categoriesRepo := repository.NewCategoriesRepo(db)
	categoriesUseCase := usecase.NewCategoriesUseCase(categoriesRepo, logger)
	statisticsRepo := repository.NewStatisticsRepo(db)
	statisticsUseCase := usecase.NewStatisticsUseCase(statisticsRepo)
	authRepo := repository.NewAuthRepo()
	authUseCase := usecase.NewAuthUseCase(authRepo)
	categoriesManagerUseCase := usecase.NewCategoriesManager(categoriesRepo, logger)
	repository.Migration(db, logger)

	gin.SetMode(gin.TestMode)
	handler := gin.New()
	NewRouter(
		handler,
		dreamsUseCase,
		privateDreamsUseCase,
		categoriesUseCase,
		statisticsUseCase,
		categoriesUseCase,
		authUseCase,
		categoriesManagerUseCase,
	)
	return handler
}

type apiTest struct {
	name       string
	method     string
	url        string
	body       any
	statusCode int
	response   any
	token      func() *string
	after      func(resp *httptest.ResponseRecorder)
}

func (test apiTest) evaluate(t *testing.T, h *gin.Engine) {
	t.Run(test.name, func(t *testing.T) {
		var body io.Reader
		if test.body != nil {
			marBody, _ := json.Marshal(test.body)
			body = bytes.NewBuffer(marBody)
		}
		req, _ := http.NewRequest(test.method, test.url, body)
		if test.token != nil {
			if test.token() != nil {
				req.Header.Set("Authorization", *test.token())
			}
		}
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		assert.Equal(t, test.statusCode, resp.Code)

		if test.response != nil {
			expected := test.response
			actual := reflect.New(reflect.TypeOf(expected)).Interface()
			err := json.Unmarshal(resp.Body.Bytes(), &actual)
			assert.NoError(t, err)
			assert.Equal(t, expected, reflect.Indirect(reflect.ValueOf(actual)).Interface())
		}
		if test.after != nil {
			test.after(resp)
		}
	})
}
