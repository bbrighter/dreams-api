package controller

import (
	"bytes"
	"io"
	"net/http"

	dt "github.com/bbrighter/dreams-api/internal/dreamTranscriptions"
	"github.com/bbrighter/dreams-api/internal/entities"
	"github.com/bbrighter/dreams-api/internal/recordings"
	"github.com/gin-gonic/gin"
)

type recordingsRoute struct {
	dts      *dt.DreamTranscriptionService
	rec      *recordings.RecordingsService
	mimeType string
}

func newRecordingRoute(
	handler *gin.RouterGroup,
	dts *dt.DreamTranscriptionService,
	rec *recordings.RecordingsService,
	mimeType string,
) {
	route := &recordingsRoute{dts: dts, rec: rec, mimeType: mimeType}

	dreamsGroup := handler.Group("/dreams/:dreamId")
	{
		dreamsGroup.POST("/recording", route.PostRecording)
		dreamsGroup.POST("/transcribe", route.PostTranscribe)
	}
	recordingsGroup := handler.Group("/recordings/:id")
	{
		recordingsGroup.DELETE("", route.DeleteRecording)
		recordingsGroup.GET("", route.GetRecording)
	}
}

// @Description	Transcribe a recording
// @Success		200	{object}	RecordingResponse
// @Failure		400	{object}	ErrorResponse
// @Failure		404	{object}	ErrorResponse
// @Router			/dreams/{dreamId}/transcribe [post]
// @Produce		json
func (r *recordingsRoute) PostTranscribe(g *gin.Context) {
	id, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}
	ctx := g.Request.Context()

	recording, err := r.dts.Transcribe(ctx, id)
	if handleError(g, err) {
		return
	}

	g.JSON(http.StatusOK, recordToResponse(recording))
}

type RecordingResponse struct {
	ID                 uint                        `json:"id"`
	Transcript         string                      `json:"transcript"`
	RecordingState     entities.RecordingState     `json:"recordingState"`
	TranscriptionState entities.TranscriptionState `json:"transcriptionState"`
	MimeType           string                      `json:"mimeType"`
}

func recordToResponse(rec *entities.Recording) *RecordingResponse {
	if rec == nil {
		return nil
	}

	return &RecordingResponse{
		ID:                 rec.ID,
		Transcript:         rec.Transcript,
		RecordingState:     rec.RecordingState,
		TranscriptionState: rec.TranscriptionState,
		MimeType:           rec.MimeType,
	}
}

// @Description	Add a recording to a dream
// @Success		200	{object}	RecordingResponse
// @Accept			multipart/form-data
// @Param			file	formData	file	true	"Audio recording"
// @Failure		400		{object}	ErrorResponse
// @Failure		404		{object}	ErrorResponse
// @Router			/dreams/{dreamId}/recording [post]
// @Produce		json
func (r *recordingsRoute) PostRecording(g *gin.Context) {
	id, err := parseParamUint(g, "dreamId")
	if handleError(g, err) {
		return
	}
	ctx := g.Request.Context()

	fileHeader, err := g.FormFile("file")
	if handleError(g, err) {
		return
	}
	file, err := fileHeader.Open()
	if handleError(g, err) {
		return
	}
	defer file.Close()

	blob, err := io.ReadAll(file)
	if handleError(g, err) {
		return
	}

	recording, err := r.rec.Create(ctx, id, blob, fileHeader.Header.Get("Content-Type"))
	if handleError(g, err) {
		return
	}
	g.JSON(http.StatusOK, recordToResponse(recording))
}

// @Description	Delete a recording
// @Success		200
// @Failure		400	{object}	ErrorResponse
// @Failure		404	{object}	ErrorResponse
// @Router			/recordings/{id} [delete]
func (r *recordingsRoute) DeleteRecording(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if handleError(g, err) {
		return
	}
	ctx := g.Request.Context()

	err = r.rec.Delete(ctx, id)
	if handleError(g, err) {
		return
	}

	g.Status(http.StatusOK)
}

// @Description	Get a recording
// @Produce		octet-stream
// @Success		200	{file}		binary	"Recording file"
// @Failure		400	{object}	ErrorResponse
// @Failure		404	{object}	ErrorResponse
// @Router			/recordings/{id} [get]
func (r *recordingsRoute) GetRecording(g *gin.Context) {
	id, err := parseParamUint(g, "id")
	if handleError(g, err) {
		return
	}
	ctx := g.Request.Context()

	file, err := r.rec.Get(ctx, id)
	if handleError(g, err) {
		return
	}
	reader := bytes.NewReader(file)

	g.DataFromReader(http.StatusOK, reader.Size(), r.mimeType, reader, nil)
}
