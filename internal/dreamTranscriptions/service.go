package dreamtranscriptions

import (
	"context"
	"errors"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/entities"
)

type TranscribeManager interface {
	Init() error
	Transcribe(context.Context, []byte) (string, error)
}

type DreamManager interface {
	GetDreamById(context.Context, uint) (entities.Dream, error)
	UpdateDream(context.Context, uint, dreams.UpdateDreamParams) error
}

type RecordingManager interface {
	Get(context.Context, uint) ([]byte, error)
	StartProcessing(context.Context, *entities.Recording) error
	FailProcessing(context.Context, *entities.Recording) error
	CompleteProcessing(context.Context, *entities.Recording, string) error
}

type DreamTranscriptionService struct {
	i TranscribeManager
	d DreamManager
	r RecordingManager
}

func NewDreamTranscriptionService(
	i TranscribeManager,
	d DreamManager,
	r RecordingManager,
) *DreamTranscriptionService {
	return &DreamTranscriptionService{i: i, d: d, r: r}
}

func (s *DreamTranscriptionService) Init() error {
	return s.i.Init()
}

func (s *DreamTranscriptionService) Transcribe(ctx context.Context, dreamId uint) (*entities.Recording, error) {
	dream, err := s.d.GetDreamById(ctx, dreamId)
	if err != nil {
		return &entities.Recording{}, err
	}

	rec := dream.Recording
	if dream.Recording == nil {
		return rec, apperrors.ErrNotFound
	}

	if err := s.r.StartProcessing(ctx, rec); err != nil {
		return rec, err
	}

	fail := func(err error) (*entities.Recording, error) {
		rec.TranscriptionState = entities.TranscriptionStateFailed
		return rec, errors.Join(err, s.r.FailProcessing(ctx, rec))
	}

	file, err := s.r.Get(ctx, rec.ID)
	if err != nil {
		return fail(err)
	}

	transcript, err := s.i.Transcribe(ctx, file)
	if err != nil {
		return fail(err)
	}

	err = s.r.CompleteProcessing(ctx, rec, transcript)
	if err != nil {
		return fail(err)
	}
	rec.TranscriptionState = entities.TranscriptionStateCompleted

	return rec, nil
}
