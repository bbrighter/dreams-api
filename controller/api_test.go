package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"testing"

	"github.com/bbrighter/dreams-api/internal/blob"
	"github.com/bbrighter/dreams-api/internal/categories"
	dreamcategories "github.com/bbrighter/dreams-api/internal/dreamCategories"
	dreamtranscriptions "github.com/bbrighter/dreams-api/internal/dreamTranscriptions"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/migrations"
	"github.com/bbrighter/dreams-api/internal/recordings"
	"github.com/bbrighter/dreams-api/internal/statistics"
	"github.com/bbrighter/dreams-api/internal/transcription"
	unitofwork "github.com/bbrighter/dreams-api/internal/unitOfWork"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type apiTestSuite struct {
	suite.Suite
	g        *gin.Engine
	db       *gorm.DB
	trans    *mockTrans
	mimeType string
}

type mockTrans struct{ mock.Mock }

func (m *mockTrans) ModelExists(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *mockTrans) DownloadModel(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockTrans) Transcribe(ctx context.Context, file []byte) (string, error) {
	args := m.Called(ctx, file)
	return args.String(0), args.Error(1)
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(apiTestSuite))
}

func (s *apiTestSuite) SetupTest() {
	logger := zap.NewExample()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)
	err = db.Exec("PRAGMA foreign_keys = ON").Error
	s.Require().NoError(err)
	s.db = db
	s.mimeType = "audio/webm;codecs=opus"

	migrations.Migration(db, logger)
	dreamsR := dreams.NewDreamsRepo(db)
	cats := categories.NewCategoriesRepo(db)
	stats := statistics.NewStatisticsRepo(db)
	rec := recordings.NewRecordingsRepository(db)
	s.trans = new(mockTrans)
	uow := unitofwork.NewUnitOfWork(db)
	blob := blob.NewBlobStorage("test_data")

	dreamService := dreams.NewDreamsService(dreamsR)
	recService := recordings.NewRecordingsService(rec, blob, s.mimeType)
	transService := transcription.NewTranscriptionService(s.trans)

	gin.SetMode(gin.TestMode)
	handler := gin.New()
	NewRouter(
		handler,
		dreamService,
		categories.NewCategoriesService(cats),
		dreamcategories.NewDreamCategoriesService(dreamsR, cats, uow),
		statistics.NewStatisticsService(stats),
		dreamtranscriptions.NewDreamTranscriptionService(transService, dreamService, recService),
		recService,
		"mp3",
	)

	s.g = handler
}

type apiTest struct {
	name        string
	method      string
	url         string
	body        any
	bodyReader  io.Reader
	statusCode  int
	response    any
	contentType string
}

func (s *apiTestSuite) evaluate(test apiTest) {
	s.Run(test.name, func() {
		var body io.Reader

		if test.bodyReader != nil {
			body = test.bodyReader
		} else if test.body != nil {
			marBody, err := json.Marshal(test.body)
			s.Require().NoError(err)
			body = bytes.NewBuffer(marBody)
		}

		req, err := http.NewRequest(test.method, test.url, body)
		s.Require().NoError(err)

		if test.contentType != "" {
			req.Header.Set("Content-Type", test.contentType)
		} else if test.body != nil {
			req.Header.Set("Content-Type", "application/json")
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
	})
}

func recordingMultipartBody(
	t *testing.T,
	filename string,
	content []byte,
	mimeType string,
) (io.Reader, string) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	header := make(textproto.MIMEHeader)
	header.Add("Content-Type", mimeType)
	header.Set("Content-Disposition", multipart.FileContentDisposition("file", filename))
	part, err := writer.CreatePart(header)

	require.NoError(t, err)

	_, err = part.Write(content)
	require.NoError(t, err)

	require.NoError(t, writer.Close())

	return &body, writer.FormDataContentType()
}
