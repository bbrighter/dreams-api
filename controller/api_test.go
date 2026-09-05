package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bbrighter/dreams-api/internal/categories"
	dreamcategories "github.com/bbrighter/dreams-api/internal/dreamCategories"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/migrations"
	"github.com/bbrighter/dreams-api/internal/statistics"
	unitofwork "github.com/bbrighter/dreams-api/internal/unitOfWork"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApiTestSuite struct {
	suite.Suite
	g  *gin.Engine
	db *gorm.DB
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (s *ApiTestSuite) SetupTest() {
	logger := zap.NewExample()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)
	err = db.Exec("PRAGMA foreign_keys = ON").Error
	s.Require().NoError(err)
	s.db = db

	migrations.Migration(db, logger)
	dreamsR := dreams.NewDreamsRepo(db)
	cats := categories.NewCategoriesRepo(db)
	stats := statistics.NewStatisticsRepo(db)
	uow := unitofwork.NewUnitOfWork(db)

	gin.SetMode(gin.TestMode)
	handler := gin.New()
	NewRouter(
		handler,
		dreams.NewDreamsService(dreamsR),
		categories.NewCategoriesService(cats),
		dreamcategories.NewDreamCategoriesService(dreamsR, cats, uow),
		statistics.NewStatisticsService(stats),
	)

	s.g = handler
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

func (s *ApiTestSuite) evaluate(test apiTest) {
	s.Run(test.name, func() {
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
		s.g.ServeHTTP(resp, req)
		s.Equal(test.statusCode, resp.Code)

		if test.response != nil {
			expected := test.response
			actual := reflect.New(reflect.TypeOf(expected)).Interface()
			err := json.Unmarshal(resp.Body.Bytes(), &actual)
			s.NoError(err)
			s.Equal(expected, reflect.Indirect(reflect.ValueOf(actual)).Interface())
		}
		if test.after != nil {
			test.after(resp)
		}
	})
}
