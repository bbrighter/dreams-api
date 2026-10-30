package controller

import (
	"net/http"
	"time"

	"github.com/bbrighter/dreams-api/internal/entities"
	"github.com/stretchr/testify/mock"
)

func (s *apiTestSuite) TestPostTranscribe() {
	s.trans.On("Transcribe", mock.Anything, mock.Anything).Return("text", nil)

	content := []byte("data")
	body, contentType := recordingMultipartBody(s.T(), "recording.webm", content, s.mimeType)

	tests := []apiTest{
		{
			name: "add dream", method: http.MethodPost, url: "/dreams", body: PostDreamRequest{Date: time.Date(2026, 9, 5, 13, 10, 0, 0, time.UTC)},
			statusCode: http.StatusCreated,
		},
		{
			name: "add a recording", method: http.MethodPost, url: "/dreams/1/recording", bodyReader: body, contentType: contentType,
			response: RecordingResponse{ID: 1, Transcript: "", RecordingState: entities.RecordingStateStored, TranscriptionState: entities.TranscriptionStatePending, MimeType: "audio/webm;codecs=opus"}, statusCode: http.StatusOK,
		},
		{
			name: "get recording", method: http.MethodGet, url: "/recordings/1",
			statusCode: http.StatusOK,
		},
		{
			name: "transcribe", method: http.MethodPost, url: "/dreams/1/transcribe",
			statusCode: http.StatusOK,
			response:   RecordingResponse{ID: 1, Transcript: "text", RecordingState: entities.RecordingStateStored, TranscriptionState: entities.TranscriptionStateCompleted, MimeType: "audio/webm;codecs=opus"},
		},
		{
			name: "Delete recording", method: http.MethodDelete, url: "/recordings/1",
			statusCode: http.StatusOK,
		},
		{
			name: "Recording not found", method: http.MethodGet, url: "/recordings/1",
			statusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		s.evaluate(test)
	}

}
