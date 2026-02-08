package repository

import (
	"context"
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepoTestSuite struct {
	suite.Suite
	db     *gorm.DB
	ctx    context.Context
	stats  *StatisticsRepo
	dreams *DreamsRepo
	cats   *CategoriesRepo
	mgmt   *ManagementRepo
}

func (s *RepoTestSuite) SetupSuite() {

	logger, _ := zap.NewDevelopment()
	s.db = NewDatabase(":memory:", logger)
	s.ctx = context.Background()
	Migration(s.db, logger)

	s.stats = NewStatisticsRepo(s.db)
	s.dreams = NewDreamsRepo(s.db)
	s.cats = NewCategoriesRepo(s.db)
	s.mgmt = NewManagementRepo(s.db)
}

func (s *RepoTestSuite) TearDownTest() {
	tables := []string{"categories_dreams", "dreams", "categories"}
	for _, t := range tables {
		err := s.db.Exec("DELETE FROM " + t).Error
		s.Require().NoError(err)
	}
}

func (s *RepoTestSuite) TearDownSubTest() {
	tables := []string{"categories_dreams", "dreams", "categories"}
	for _, t := range tables {
		err := s.db.Exec("DELETE FROM " + t).Error
		s.Require().NoError(err)
	}
	// var err error
	// err = s.db.Exec("DELETE FROM categories").Error
	// s.Require().NoError(err)
	// err = s.db.Exec("DELETE FROM categories_dreams").Error
	// s.Require().NoError(err)
	// err = s.db.Exec("DELETE FROM dreams").Error
	// s.Require().NoError(err)
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

func (s *RepoTestSuite) CreateDream(cat bool, person bool) *entity.Dream {
	var dream = &entity.Dream{Date: time.Date(2022, 2, 1, 2, 2, 2, 2, time.UTC), Categories: entity.Categories{}}
	if cat {
		var cat1 = entity.Category{Name: "CatName", Type: entity.TypeCategory}
		dream.Categories = append(dream.Categories, cat1)
	}
	if person {
		var per1 = entity.Category{Name: "PerName", Type: entity.TypePerson}
		dream.Categories = append(dream.Categories, per1)
	}
	err := s.db.Create(dream).Error
	s.Require().NoError(err)
	return dream
}
