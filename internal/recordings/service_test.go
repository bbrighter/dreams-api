package recordings_test

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
	"github.com/bbrighter/dreams-api/internal/entities"
	"github.com/bbrighter/dreams-api/internal/recordings"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type mockRecordRepo struct {
	mock.Mock
}

func (m *mockRecordRepo) Create(ctx context.Context, rec *entities.Recording) error {
	args := m.Called(ctx, rec)
	return args.Error(0)
}
func (m *mockRecordRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockRecordRepo) Update(ctx context.Context, id uint, values map[string]any) error {
	args := m.Called(ctx, id, values)
	return args.Error(0)
}

func (m *mockRecordRepo) First(ctx context.Context, id uint) (entities.Recording, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Recording), args.Error(1)
}

func (s *serviceTestSuite) mockCreateRecord(recId uint, err error) {
	s.rec.On("Create", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			rec := args.Get(1).(*entities.Recording)
			rec.ID = recId
		}).
		Return(err)
}
func (s *serviceTestSuite) mockUpdateRecord(err error) {
	s.rec.On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(err)
}
func (s *serviceTestSuite) mockDeleteRecord(err error) {
	s.rec.On("Delete", mock.Anything, mock.Anything).
		Return(err)
}

type mockBlobStorage struct {
	mock.Mock
}

func (m *mockBlobStorage) Create(id uint, file []byte) error {
	args := m.Called(id, file)
	return args.Error(0)
}
func (m *mockBlobStorage) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockBlobStorage) Get(id uint) ([]byte, error) {
	args := m.Called(id)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *serviceTestSuite) mockBlobCreate(err error) {
	s.blob.On("Create", mock.Anything, mock.Anything).
		Return(err)
}

func (s *serviceTestSuite) mockBlobDelete(err error) {
	s.blob.On("Delete", mock.Anything).
		Return(err)
}

type serviceTestSuite struct {
	suite.Suite
	ctx     context.Context
	rec     *mockRecordRepo
	blob    *mockBlobStorage
	service *recordings.RecordingsService
	dreamId uint
}

func (s *serviceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.rec = new(mockRecordRepo)
	s.blob = new(mockBlobStorage)

	s.service = recordings.NewRecordingsService(s.rec, s.blob, "mp3")

	s.blob.On("GetFilePath", mock.Anything).
		Return("path")
}

func TestRecordService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (s *serviceTestSuite) TestCreateOk() {
	var recordId uint = 20

	s.rec.On(
		"Create",
		s.ctx,
		// Must be matched here; r is a pointer, hence it references the wrong update
		mock.MatchedBy(func(r *entities.Recording) bool {
			return r.DreamID == 1 &&
				r.RecordingState == entities.RecordingStateCreating &&
				r.TranscriptionState == entities.TranscriptionStatePending
		}),
	).Run(func(args mock.Arguments) {
		args.Get(1).(*entities.Recording).ID = recordId
	}).Return(nil).Once()
	s.mockUpdateRecord(nil)
	s.mockBlobCreate(nil)
	s.mockBlobDelete(nil)

	rec, err := s.service.Create(s.ctx, 1, []byte{81}, "mp3")

	s.NoError(err)
	s.Equal(recordId, rec.ID)

	s.blob.AssertCalled(
		s.T(),
		"Create",
		recordId,
		[]byte{81},
	)

	s.rec.AssertCalled(
		s.T(),
		"Update",
		s.ctx,
		recordId,
		mock.MatchedBy(func(v map[string]any) bool {
			return v["recording_state"] == entities.RecordingStateStored
		}),
	)
}

func (s *serviceTestSuite) TestCreateDatabaseFails() {
	s.mockCreateRecord(0, gorm.ErrInvalidDB)

	rec, err := s.service.Create(s.ctx, 20, []byte{81}, "mp3")

	s.Error(err)
	s.ErrorIs(err, gorm.ErrInvalidDB)
	s.Zero(rec.ID)
}

func (s *serviceTestSuite) TestCreateBlobFails() {
	var recordId uint = 20
	s.mockCreateRecord(recordId, nil)
	s.mockBlobCreate(errors.ErrUnsupported)
	s.mockDeleteRecord(nil)
	s.mockBlobDelete(nil)

	rec, err := s.service.Create(s.ctx, 1, []byte{81}, "mp3")

	s.Error(err)
	s.ErrorIs(err, errors.ErrUnsupported)
	s.Equal(recordId, rec.ID)

	s.rec.AssertCalled(s.T(), "Delete", s.ctx, recordId)
	s.blob.AssertCalled(s.T(), "Delete", recordId)
}

func (s *serviceTestSuite) TestUpdateRecordFails() {
	var recordId uint = 20
	s.mockCreateRecord(recordId, nil)
	s.mockBlobCreate(nil)
	s.mockUpdateRecord(errors.ErrUnsupported)
	s.mockDeleteRecord(nil)
	s.mockBlobDelete(nil)

	rec, err := s.service.Create(s.ctx, 1, []byte{81}, "mp3")

	s.Error(err)
	s.ErrorIs(err, errors.ErrUnsupported)
	s.Equal(recordId, rec.ID)

	s.rec.AssertCalled(s.T(), "Delete", s.ctx, recordId)
	s.blob.AssertCalled(s.T(), "Delete", recordId)
}

func (s *serviceTestSuite) TestCreateBlobFailsAndDeletionAsWell() {
	var recordId uint = 20
	s.mockCreateRecord(recordId, nil)
	s.mockBlobCreate(errors.ErrUnsupported)
	s.mockDeleteRecord(gorm.ErrInvalidDB)
	s.mockBlobDelete(nil)

	rec, err := s.service.Create(s.ctx, 1, []byte{81}, "mp3")

	s.Error(err)
	s.ErrorIs(err, errors.ErrUnsupported)
	s.ErrorIs(err, gorm.ErrInvalidDB)
	s.Equal(recordId, rec.ID)

	s.rec.AssertCalled(s.T(), "Delete", s.ctx, recordId)
	s.blob.AssertCalled(s.T(), "Delete", recordId)
}

func (s *serviceTestSuite) TestCreateBadMimeType() {
	_, err := s.service.Create(s.ctx, 1, []byte{81}, "bad mime")

	s.Error(err)
	s.ErrorIs(err, apperrors.ErrBadParam)
}
