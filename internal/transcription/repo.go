package transcription

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
)

type TranscriptionRepo struct {
	client   *http.Client
	url      string
	model    string
	language string
	mimeType string
}

func NewTranscriptionRepo(url string, model string, mimeType string) *TranscriptionRepo {
	return &TranscriptionRepo{
		url:      url,
		model:    model,
		language: "de",
		mimeType: mimeType,
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (r *TranscriptionRepo) ModelExists(ctx context.Context) (bool, error) {
	url := r.url + "/v1/models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, apperrors.ErrWhisper
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	type modelResp struct {
		Data []struct {
			Id string `json:"id"`
		} `json:"data"`
	}
	var model modelResp

	if err := json.Unmarshal(body, &model); err != nil {
		return false, err
	}
	if model.Data == nil {
		return false, apperrors.ErrBadParam
	}

	for _, d := range model.Data {
		if d.Id == r.model {
			return true, nil
		}
	}
	return false, nil
}

func (r *TranscriptionRepo) DownloadModel(ctx context.Context) error {
	url := r.url + "/v1/models/" + r.model

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return apperrors.ErrWhisper
	}

	return nil
}

func (r *TranscriptionRepo) Transcribe(ctx context.Context, file []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("model", r.model); err != nil {
		return "", err
	}
	if err := writer.WriteField("language", r.language); err != nil {
		return "", err
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="audio.wav"`)
	header.Set("Content-Type", r.mimeType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(file); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	url := r.url + "/v1/audio/transcriptions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", apperrors.ErrWhisper
	}

	type transcriptionResponse struct {
		Text string `json:"text"`
	}
	var result transcriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Text, nil

}
