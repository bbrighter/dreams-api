package transcription

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelExists(t *testing.T) {
	tests := map[string]struct {
		statusCode    int
		body          string
		expectedErr   error
		expectedFound bool
	}{
		"found": {
			statusCode:    http.StatusOK,
			body:          `{"data": [{"id": "modelName", "created": 12345}]}`,
			expectedFound: true,
		},
		"found in multiple models": {
			statusCode:    http.StatusOK,
			body:          `{"data": [{"id": "wrongName"},{"id": "modelName"}]}`,
			expectedFound: true,
		},
		"not found": {
			statusCode:    http.StatusOK,
			body:          `{"data": [{"id": "wrong", "created": 12345}]}`,
			expectedFound: false,
		},
		"http error": {
			statusCode:    http.StatusBadRequest,
			body:          "",
			expectedFound: false,
			expectedErr:   apperrors.ErrWhisper,
		},
		"bad format": {
			statusCode:    http.StatusOK,
			body:          `{"strange": 123}`,
			expectedErr:   apperrors.ErrBadParam,
			expectedFound: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "/v1/models", r.URL.Path)

					w.WriteHeader(test.statusCode)
					w.Write([]byte(test.body))
				},
			))
			defer server.Close()

			repo := NewTranscriptionRepo(server.URL, "modelName", "mime")

			found, err := repo.ModelExists(t.Context())
			if test.expectedErr != nil {
				assert.ErrorIs(t, err, test.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedFound, found)
		})
	}

}

func TestDownloadMode(t *testing.T) {
	tests := map[string]struct {
		modelName   string
		statusCode  int
		expectedErr error
	}{
		"ok":          {modelName: "xyz", statusCode: http.StatusOK},
		"wrong model": {modelName: "abc", statusCode: http.StatusNotFound, expectedErr: apperrors.ErrWhisper},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					assert.Equal(t, "/v1/models/xyz", r.URL.Path)

					w.WriteHeader(test.statusCode)
				},
			))
			defer server.Close()

			repo := NewTranscriptionRepo(server.URL, "xyz", "mime")

			err := repo.DownloadModel(context.Background())
			if test.expectedErr != nil {
				assert.ErrorIs(t, err, test.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestTranscribe(t *testing.T) {
	tests := map[string]struct {
		statusCode   int
		responseBody string
		expectedErr  bool
	}{
		"ok": {
			statusCode:   http.StatusOK,
			responseBody: `{"text":"hello world"}`,
		},
		"bad request": {
			statusCode:  http.StatusBadRequest,
			expectedErr: true,
		},
		"invalid json": {
			statusCode:   http.StatusOK,
			responseBody: "xyz",
			expectedErr:  true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					assert.Equal(t, "/v1/audio/transcriptions", r.URL.Path)
					assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

					err := r.ParseMultipartForm(1 << 20)
					require.NoError(t, err)

					assert.Equal(t, "modelName", r.FormValue("model"))

					assert.Equal(t, "de", r.FormValue("language"))

					audio, header, err := r.FormFile("file")
					require.NoError(t, err)
					defer audio.Close()
					assert.Equal(t, "audio.wav", header.Filename)
					contents, err := io.ReadAll(audio)
					require.NoError(t, err)
					assert.Equal(t, []byte("123"), contents)

					w.WriteHeader(test.statusCode)
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(test.responseBody))
				},
			))
			defer server.Close()

			repo := NewTranscriptionRepo(server.URL, "modelName", "mime")

			transcript, err := repo.Transcribe(t.Context(), []byte("123"))
			if test.expectedErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, "hello world", transcript)
		})
	}
}
