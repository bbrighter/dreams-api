package recordings_test

import (
	"context"
	"testing"

	"github.com/bbrighter/dreams-api/internal/entities"
	"github.com/bbrighter/dreams-api/internal/migrations"
	"github.com/bbrighter/dreams-api/internal/recordings"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type recordingsTestSuite struct {
	suite.Suite
	repo *recordings.RecordingsRepository
	db   *gorm.DB
	ctx  context.Context
}

func (s *recordingsTestSuite) SetupSuite() {
	logger, _ := zap.NewDevelopment()
	s.db = migrations.NewDatabase(":memory:", "", logger)
	s.ctx = context.Background()
	migrations.Migration(s.db, logger)

	s.repo = recordings.NewRecordingsRepository(s.db)
}

func (s *recordingsTestSuite) TearDownTest() {
	tables := []string{"recordings", "dreams"}
	for _, t := range tables {
		err := s.db.Exec("DELETE FROM " + t).Error
		s.Require().NoError(err)
	}
}

func (s *recordingsTestSuite) createTestDream() uint {
	dream := &entities.Dream{}
	err := s.db.Create(dream).Error
	s.Require().NoError(err)
	return dream.ID
}

func (s *recordingsTestSuite) createTestRecording() entities.Recording {
	dreamId := s.createTestDream()
	rec := entities.Recording{DreamID: dreamId}
	err := s.db.Create(&rec).Error
	s.Require().NoError(err)
	return rec
}

func TestRecordingTestSuite(t *testing.T) {
	suite.Run(t, new(recordingsTestSuite))
}

func (s *recordingsTestSuite) TestCreateRecording() {
	dreamId := s.createTestDream()

	rec := &entities.Recording{DreamID: dreamId}
	err := s.repo.Create(s.ctx, rec)

	s.NoError(err)
}

func (s *recordingsTestSuite) TestCreateTwoRecordingsForOneDreamNotPossible() {
	dreamId := s.createTestDream()

	rec := &entities.Recording{DreamID: dreamId}
	err := s.repo.Create(s.ctx, rec)
	s.NoError(err)

	rec2 := &entities.Recording{DreamID: dreamId}
	err2 := s.repo.Create(s.ctx, rec2)
	s.Error(err2)
	s.ErrorIs(err2, gorm.ErrDuplicatedKey)
}

func (s *recordingsTestSuite) TestDeleteRecording() {
	rec := s.createTestRecording()

	err := s.repo.Delete(s.ctx, rec.ID)

	s.NoError(err)
	c, err := gorm.G[entities.Recording](s.db).Count(s.ctx, "*")
	s.NoError(err)
	s.EqualValues(0, c)
}

func (s *recordingsTestSuite) TestDeleteDreamDeletesRecordingAsWell() {
	rec := s.createTestRecording()

	affectedRows, err := gorm.G[entities.Dream](s.db).Where("id = ?", rec.DreamID).Delete(s.ctx)
	s.NoError(err)
	s.EqualValues(1, affectedRows)

	c, err := gorm.G[entities.Recording](s.db).Count(s.ctx, "*")
	s.NoError(err)
	s.EqualValues(0, c)
}
