package dreamtranscriptions

import (
	"context"
	"testing"

	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTranscribeManager struct{ mock.Mock }

func (m *mockTranscribeManager) Init() error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockTranscribeManager) Transcribe(ctx context.Context, file []byte) (string, error) {
	args := m.Called(ctx, file)
	return args.String(0), args.Error(1)
}

type mockDreamManager struct{ mock.Mock }

func (m *mockDreamManager) GetDreamById(ctx context.Context, id uint) (entities.Dream, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Dream), args.Error(1)
}
func (m *mockDreamManager) UpdateDream(ctx context.Context, id uint, params dreams.UpdateDreamParams) error {
	args := m.Called(ctx, id, params)
	return args.Error(0)
}

type mockRecordingManager struct{ mock.Mock }

func (m *mockRecordingManager) Get(ctx context.Context, id uint) ([]byte, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *mockRecordingManager) StartProcessing(ctx context.Context, rec *entities.Recording) error {
	args := m.Called(ctx, rec)
	return args.Error(0)
}
func (m *mockRecordingManager) FailProcessing(ctx context.Context, rec *entities.Recording) error {
	args := m.Called(ctx, rec)
	return args.Error(0)
}
func (m *mockRecordingManager) CompleteProcessing(ctx context.Context, rec *entities.Recording, desc string) error {
	args := m.Called(ctx, rec, desc)
	return args.Error(0)
}

func TestTranscribe(t *testing.T) {
	tests := map[string]struct {
		name       string
		setupMocks func(
			*mockTranscribeManager,
			*mockDreamManager,
			*mockRecordingManager,
			context.CancelFunc,
			*entities.TranscriptionState,
		)
		wantErr        error
		wantTranscript string
		wantState      entities.TranscriptionState
	}{
		"successful transcription": {
			setupMocks: func(tm *mockTranscribeManager, dm *mockDreamManager, rm *mockRecordingManager, _ context.CancelFunc, persistedState *entities.TranscriptionState) {
				dm.On("GetDreamById", mock.Anything, uint(1)).
					Return(entities.Dream{
						ID: 1,
						Recording: &entities.Recording{
							ID:             1,
							RecordingState: entities.RecordingStateStored,
						},
					}, nil)

				rm.On("Get", mock.Anything, uint(1)).
					Return([]byte{}, nil)

				rm.On("StartProcessing", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						*persistedState = entities.TranscriptionStateProcessing
					}).
					Return(nil)

				tm.On("Transcribe", mock.Anything, mock.Anything).
					Return("transcription", nil)

				rm.On("CompleteProcessing", mock.Anything, mock.Anything, "transcription").
					Run(func(args mock.Arguments) {
						rec := args.Get(1).(*entities.Recording)
						rec.Transcript = args.String(2)
						*persistedState = entities.TranscriptionStateCompleted
					}).
					Return(nil)
			},
			wantErr:        nil,
			wantTranscript: "transcription",
			wantState:      entities.TranscriptionStateCompleted,
		},
		"transcription times out": {
			setupMocks: func(tm *mockTranscribeManager, dm *mockDreamManager, rm *mockRecordingManager, cancel context.CancelFunc, persistedState *entities.TranscriptionState) {
				dm.On("GetDreamById", mock.Anything, uint(1)).
					Return(entities.Dream{
						ID: 1,
						Recording: &entities.Recording{
							ID:             1,
							RecordingState: entities.RecordingStateStored,
						},
					}, nil)

				rm.On("Get", mock.Anything, uint(1)).
					Return([]byte{}, nil)

				rm.On("StartProcessing", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						*persistedState = entities.TranscriptionStateProcessing
					}).
					Return(nil)

				tm.On("Transcribe", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						cancel()
					}).
					Return("", context.Canceled)

				rm.On("FailProcessing", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						cleanupCtx := args.Get(0).(context.Context)
						assert.NoError(t, cleanupCtx.Err())
						if cleanupCtx.Err() == nil {
							*persistedState = entities.TranscriptionStateFailed
						}
					}).
					Return(nil)
			},
			wantErr:   context.Canceled,
			wantState: entities.TranscriptionStateFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := new(mockTranscribeManager)
			dm := new(mockDreamManager)
			rm := new(mockRecordingManager)
			persistedState := entities.TranscriptionStatePending
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			tt.setupMocks(tm, dm, rm, cancel, &persistedState)

			service := NewDreamTranscriptionService(tm, dm, rm)

			rec, err := service.Transcribe(ctx, 1)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.wantState, rec.TranscriptionState)
				rm.AssertCalled(t, "FailProcessing", mock.Anything, mock.Anything)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTranscript, rec.Transcript)
				assert.Equal(t, tt.wantState, rec.TranscriptionState)
			}
			assert.Equal(t, tt.wantState, persistedState)

			tm.AssertExpectations(t)
			dm.AssertExpectations(t)
			rm.AssertExpectations(t)
		})
	}
}
