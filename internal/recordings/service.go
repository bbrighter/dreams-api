package recordings

import (
	"context"
	"errors"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
	"github.com/bbrighter/dreams-api/internal/entities"
)

type RecordManager interface {
	Create(context.Context, *entities.Recording) error
	Delete(context.Context, uint) error
	Update(context.Context, uint, map[string]any) error
	First(context.Context, uint) (entities.Recording, error)
}

type BlobManager interface {
	Create(uint, []byte) error
	Delete(uint) error
	Get(uint) ([]byte, error)
}

type RecordingsService struct {
	record   RecordManager
	blob     BlobManager
	mimeType string
}

func NewRecordingsService(r RecordManager, b BlobManager, mimeType string) *RecordingsService {
	return &RecordingsService{record: r, blob: b, mimeType: mimeType}
}

func (s *RecordingsService) Create(ctx context.Context, dreamId uint, file []byte, mimeType string) (*entities.Recording, error) {
	rec := &entities.Recording{
		DreamID:            dreamId,
		RecordingState:     entities.RecordingStateCreating,
		TranscriptionState: entities.TranscriptionStatePending,
		MimeType:           mimeType,
	}
	if mimeType != s.mimeType {
		return rec, apperrors.ErrBadParam
	}

	if err := s.record.Create(ctx, rec); err != nil {
		return rec, err
	}

	if err := s.blob.Create(rec.ID, file); err != nil {
		return rec, errors.Join(
			err,
			s.record.Delete(ctx, rec.ID),
			s.blob.Delete(rec.ID),
		)

	}

	if err := s.record.Update(ctx, rec.ID, map[string]any{
		"recording_state": entities.RecordingStateStored,
	}); err != nil {
		return rec, errors.Join(
			err,
			s.record.Delete(ctx, rec.ID),
			s.blob.Delete(rec.ID),
		)
	}
	rec.RecordingState = entities.RecordingStateStored

	return rec, nil
}

func (s *RecordingsService) Get(ctx context.Context, id uint) ([]byte, error) {
	rec, err := s.record.First(ctx, id)
	if err != nil {
		return []byte{}, err
	}
	if rec.RecordingState != entities.RecordingStateStored {
		return []byte{}, apperrors.ErrWrongState
	}
	return s.blob.Get(id)
}

func (s *RecordingsService) AdaptRecordState(ctx context.Context, rec *entities.Recording, state entities.TranscriptionState) error {
	rec.TranscriptionState = state
	return s.record.Update(
		ctx,
		rec.ID,
		map[string]any{
			"transcription_state": state,
		},
	)
}

func (s *RecordingsService) StartProcessing(ctx context.Context, rec *entities.Recording) error {
	return s.AdaptRecordState(ctx, rec, entities.TranscriptionStateProcessing)
}

func (s *RecordingsService) FailProcessing(ctx context.Context, rec *entities.Recording) error {
	return s.AdaptRecordState(ctx, rec, entities.TranscriptionStateFailed)
}

func (s *RecordingsService) CompleteProcessing(ctx context.Context, rec *entities.Recording, transcript string) error {
	rec.TranscriptionState = entities.TranscriptionStateCompleted
	rec.Transcript = transcript
	return s.record.Update(
		ctx,
		rec.ID,
		map[string]any{
			"transcription_state": entities.TranscriptionStateCompleted,
			"transcript":          transcript,
		},
	)
}

func (s *RecordingsService) Delete(ctx context.Context, id uint) error {
	if err := s.blob.Delete(id); err != nil {
		return err
	}
	return s.record.Delete(ctx, id)
}
