package transcription

import "context"

type TransactionManager interface {
	ModelExists(context.Context) (bool, error)
	DownloadModel(context.Context) error
	Transcribe(context.Context, []byte) (string, error)
}

type TranscriptionService struct {
	r TransactionManager
}

func NewTranscriptionService(r TransactionManager) *TranscriptionService {
	return &TranscriptionService{r: r}
}

func (s *TranscriptionService) Init() error {
	ctx := context.Background()
	exists, err := s.r.ModelExists(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.r.DownloadModel(ctx)
}

func (s *TranscriptionService) Transcribe(ctx context.Context, file []byte) (string, error) {
	return s.r.Transcribe(ctx, file)
}
