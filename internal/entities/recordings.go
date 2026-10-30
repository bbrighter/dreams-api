package entities

type RecordingState string

const (
	RecordingStateCreating RecordingState = "creating"
	RecordingStateStored   RecordingState = "stored"
)

type TranscriptionState string

const (
	TranscriptionStatePending    TranscriptionState = "pending"
	TranscriptionStateProcessing TranscriptionState = "processing"
	TranscriptionStateCompleted  TranscriptionState = "completed"
	TranscriptionStateFailed     TranscriptionState = "failed"
)

type Recording struct {
	ID                 uint
	Dream              Dream `gorm:"constraint:OnDelete:CASCADE;"`
	DreamID            uint  `gorm:"unique"`
	RecordingState     RecordingState
	TranscriptionState TranscriptionState
	Transcript         string
	MimeType           string
}
