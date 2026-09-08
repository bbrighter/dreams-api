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

func (s *repoTestSuite) TestCountByCategoryAndMonth_noData() {
	stats, err := s.repo.CountByCategoryAndMonth(s.ctx)

	s.NoError(err)
	s.Len(stats, 0)
}

func (s *repoTestSuite) TestCountByCategoryAndMonth_twoDreamsWithEachOneCategory() {
	var cat1Id uint = 1
	var cat2Id uint = 2
	february2026 := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	s.createDreamWith(february2026, []entities.Category{{ID: cat1Id, Name: "Cat", Type: entities.TypeCategory}})
	s.createDreamWith(february2026, []entities.Category{{ID: cat2Id, Name: "Person", Type: entities.TypePerson}})

	stats, err := s.repo.CountByCategoryAndMonth(s.ctx)

	s.NoError(err)
	s.Len(stats, 3)

	s.Contains(stats, CountByCatAndMonth{ResultType: "category", CategoryId: &cat1Id, Month: "2026/02", Count: 1})
	s.Contains(stats, CountByCatAndMonth{ResultType: "category", CategoryId: &cat2Id, Month: "2026/02", Count: 1})
	s.Contains(stats, CountByCatAndMonth{ResultType: "total", CategoryId: nil, Month: "2026/02", Count: 2})
}

func (s *repoTestSuite) TestCountByCategoryAndMonth_oneDreamWithTwoCategories() {
	var cat1Id uint = 1
	var cat2Id uint = 2
	february2026 := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	s.createDreamWith(february2026, []entities.Category{
		{ID: cat1Id, Name: "Cat", Type: entities.TypeCategory},
		{ID: cat2Id, Name: "Person", Type: entities.TypePerson},
	})

	stats, err := s.repo.CountByCategoryAndMonth(s.ctx)

	s.NoError(err)
	s.Len(stats, 3)

	s.Contains(stats, CountByCatAndMonth{ResultType: "category", CategoryId: &cat1Id, Month: "2026/02", Count: 1})
	s.Contains(stats, CountByCatAndMonth{ResultType: "category", CategoryId: &cat2Id, Month: "2026/02", Count: 1})
	s.Contains(stats, CountByCatAndMonth{ResultType: "total", CategoryId: nil, Month: "2026/02", Count: 1})
}

func (s *repoTestSuite) TestCountByCategoryAndMonth_multipleMonths() {
	var cat1Id uint = 1
	var cat2Id uint = 2

	february2026 := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	march2026 := february2026.Add(time.Hour * 24 * 31)
	s.createDreamWith(february2026, []entities.Category{{ID: cat1Id, Name: "Cat", Type: entities.TypeCategory}})
	s.createDreamWith(february2026, []entities.Category{{ID: cat2Id, Name: "Person", Type: entities.TypePerson}})
	s.createDreamWith(march2026, []entities.Category{{ID: cat2Id, Name: "Person", Type: entities.TypePerson}})

	stats, err := s.repo.CountByCategoryAndMonth(s.ctx)

	s.NoError(err)
	s.Len(stats, 5)

	s.Contains(stats, CountByCatAndMonth{ResultType: "total", CategoryId: nil, Month: "2026/02", Count: 2})
	s.Contains(stats, CountByCatAndMonth{ResultType: "category", CategoryId: &cat1Id, Month: "2026/02", Count: 1})
	s.Contains(stats, CountByCatAndMonth{ResultType: "category", CategoryId: &cat2Id, Month: "2026/02", Count: 1})

	s.Contains(stats, CountByCatAndMonth{ResultType: "total", CategoryId: nil, Month: "2026/03", Count: 1})
	s.Contains(stats, CountByCatAndMonth{ResultType: "category", CategoryId: &cat2Id, Month: "2026/03", Count: 1})
}
