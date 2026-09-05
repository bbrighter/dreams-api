package statistics

import (
	"context"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entities"
	"github.com/bbrighter/dreams-api/internal/migrations"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	ctx  context.Context
	repo *StatisticsRepo
}

func (s *repoTestSuite) SetupSuite() {
	logger, _ := zap.NewDevelopment()
	s.db = migrations.NewDatabase(":memory:", logger)
	s.ctx = context.Background()
	s.repo = NewStatisticsRepo(s.db)
	s.db.AutoMigrate(
		entities.Category{},
		entities.Dream{},
	)
}

func (s *repoTestSuite) TearDownTest() {
	tables := []string{"categories_dreams", "dreams", "categories"}
	for _, t := range tables {
		err := s.db.Exec("DELETE FROM " + t).Error
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) TearDownSubTest() {
	tables := []string{"categories_dreams", "dreams", "categories"}
	for _, t := range tables {
		err := s.db.Exec("DELETE FROM " + t).Error
		s.Require().NoError(err)
	}
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

func (s *repoTestSuite) createDreamWith(date time.Time, cats []entities.Category) uint {
	dream := &entities.Dream{Date: date, Categories: cats}
	err := gorm.G[entities.Dream](s.db).Create(s.ctx, dream)
	s.Require().NoError(err)
	return dream.ID
}

func (s *repoTestSuite) TestCountByDreamAndMonth() {
	february2026 := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	s.createDreamWith(february2026, []entities.Category{{ID: 1, Name: "Cat", Type: entities.TypeCategory}})

	counts, err := s.repo.CountByDreamAndMonth(s.ctx)
	s.NoError(err)
	s.Len(counts, 1)
	s.Equal("2026/02", counts[0].Month)
	s.EqualValues(1, counts[0].Count)
}

func (s *repoTestSuite) TestCountByDreamAndMonth_multipleDreams() {
	february2026 := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	s.createDreamWith(february2026, []entities.Category{{ID: 1, Name: "Cat", Type: entities.TypeCategory}})
	s.createDreamWith(february2026, []entities.Category{{ID: 1, Name: "Cat", Type: entities.TypeCategory}})
	s.createDreamWith(february2026, []entities.Category{{ID: 2, Name: "Person", Type: entities.TypePerson}})

	counts, err := s.repo.CountByDreamAndMonth(s.ctx)

	s.NoError(err)
	s.Len(counts, 1)
	s.Equal("2026/02", counts[0].Month)
	s.EqualValues(3, counts[0].Count)
}

func (s *repoTestSuite) TestCountByDreamAndMonth_multipleMonths() {
	february2026 := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	march2026 := february2026.Add(time.Hour * 24 * 31)
	s.createDreamWith(february2026, []entities.Category{{ID: 1, Name: "Cat", Type: entities.TypeCategory}})
	s.createDreamWith(march2026, []entities.Category{{ID: 1, Name: "Cat", Type: entities.TypeCategory}})
	s.createDreamWith(february2026, []entities.Category{{ID: 2, Name: "Person", Type: entities.TypePerson}})

	counts, err := s.repo.CountByDreamAndMonth(s.ctx)

	s.NoError(err)
	s.Len(counts, 2)

	countsByMonth := make(map[string]int64)
	for _, count := range counts {
		countsByMonth[count.Month] = count.Count
	}
	s.Equal(map[string]int64{
		"2026/02": 2,
		"2026/03": 1,
	}, countsByMonth)
}

func (s *repoTestSuite) TestCountByCategoryAndMonth_noData() {
	stats, err := s.repo.CountByCategoryAndMonth(s.ctx)

	s.NoError(err)
	s.Len(stats, 0)
}

func (s *repoTestSuite) TestCountByCategoryAndMonth() {
	february2026 := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	s.createDreamWith(february2026, []entities.Category{{ID: 1, Name: "Cat", Type: entities.TypeCategory}})
	s.createDreamWith(february2026, []entities.Category{{ID: 2, Name: "Person", Type: entities.TypePerson}})

	stats, err := s.repo.CountByCategoryAndMonth(s.ctx)

	s.NoError(err)
	s.Len(stats, 2)

	countsByCatId := make(map[uint]CountByMonth)
	for _, stat := range stats {
		countsByCatId[stat.CategoryId] = stat.CountByMonth
	}
	s.Equal(map[uint]CountByMonth{
		1: {Month: "2026/02", Count: 1},
		2: {Month: "2026/02", Count: 1},
	}, countsByCatId)
}

func (s *repoTestSuite) TestCountByCategoryAndMonth_multipleMonths() {
	february2026 := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	march2026 := february2026.Add(time.Hour * 24 * 31)
	s.createDreamWith(february2026, []entities.Category{{ID: 1, Name: "Cat", Type: entities.TypeCategory}})
	s.createDreamWith(february2026, []entities.Category{{ID: 2, Name: "Person", Type: entities.TypePerson}})
	s.createDreamWith(march2026, []entities.Category{{ID: 2, Name: "Person", Type: entities.TypePerson}})

	stats, err := s.repo.CountByCategoryAndMonth(s.ctx)

	s.NoError(err)
	s.Len(stats, 3)

	s.Contains(stats, CountByCatAndMonth{CategoryId: 1, CountByMonth: CountByMonth{Month: "2026/02", Count: 1}})
	s.Contains(stats, CountByCatAndMonth{CategoryId: 2, CountByMonth: CountByMonth{Month: "2026/02", Count: 1}})
	s.Contains(stats, CountByCatAndMonth{CategoryId: 2, CountByMonth: CountByMonth{Month: "2026/03", Count: 1}})
}
