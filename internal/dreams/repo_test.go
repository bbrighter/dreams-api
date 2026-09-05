package dreams

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
	repo *DreamsRepo
}

func (s *repoTestSuite) SetupSuite() {
	logger, _ := zap.NewDevelopment()
	s.db = migrations.NewDatabase(":memory:", logger)
	s.ctx = context.Background()
	s.repo = NewDreamsRepo(s.db)
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

func (s *repoTestSuite) createDreamWith(desc string, finalized bool, rating *int, cats []entities.Category) uint {
	dream := &entities.Dream{Description: desc, Finalized: finalized, Rating: rating, Categories: cats}
	err := gorm.G[entities.Dream](s.db).Create(s.ctx, dream)
	s.Require().NoError(err)
	return dream.ID
}

func (s *repoTestSuite) createDream() uint {
	return s.createDreamWith("Description", false, nil, []entities.Category{})
}

func (s *repoTestSuite) TestListDreamsWithCategories() {
	s.createDreamWith("", false, nil, []entities.Category{
		{Name: "Name", Type: entities.TypeCategory},
	})

	dreams, err := s.repo.ListDreams(s.ctx)
	s.NoError(err)
	s.Len(dreams, 1)
	s.Len(dreams[0].Categories, 1)
	s.Equal("Name", dreams[0].Categories[0].Name)
	s.Equal(entities.TypeCategory, dreams[0].Categories[0].Type)
}

func (s *repoTestSuite) TestGetDreamWithCategories() {
	id := s.createDreamWith("desc", false, nil, []entities.Category{
		{Name: "Name", Type: entities.TypeCategory},
	})

	dream, err := s.repo.GetDream(s.ctx, id)
	s.NoError(err)
	s.Equal("desc", dream.Description)
	s.False(dream.Finalized)
	s.Nil(dream.Rating)
	s.Len(dream.Categories, 1)
}

func (s *repoTestSuite) TestCreateDream() {
	dream := &entities.Dream{Description: "Desc", Date: time.Date(2026, 9, 4, 19, 7, 0, 0, time.UTC)}

	id, err := s.repo.CreateDream(s.ctx, dream)
	s.NoError(err)
	s.Equal(dream.ID, id)
}

func (s *repoTestSuite) TestUpdateDream() {
	id := s.createDream()
	rating := 3
	date := time.Date(2026, 9, 4, 19, 9, 0, 0, time.UTC)
	updates := map[string]any{
		"date":        date,
		"description": "new description",
		"finalized":   true,
		"rating":      &rating,
	}

	err := s.repo.UpdateDream(s.ctx, id, updates)

	s.NoError(err)
	dream, _ := gorm.G[entities.Dream](s.db).First(s.ctx)
	s.Equal(id, dream.ID)
	s.Equal(date, dream.Date)
	s.Equal("new description", dream.Description)
	s.True(dream.Finalized)
	s.Equal(3, *dream.Rating)
}

func (s *repoTestSuite) TestUpdateDreamNotFound() {
	err := s.repo.UpdateDream(s.ctx, 100, make(map[string]any))

	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *repoTestSuite) TestDeleteDream() {
	id := s.createDream()

	err := s.repo.DeleteDream(s.ctx, id)

	s.NoError(err)
	count, _ := gorm.G[entities.Dream](s.db).Count(s.ctx, "*")
	s.EqualValues(0, count)
}

func (s *repoTestSuite) TestDeleteDreamNotFound() {
	err := s.repo.DeleteDream(s.ctx, 100)

	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *repoTestSuite) TestAddCategoryToDream() {
	id := s.createDream()
	cat := &entities.Category{Name: "Cat", Type: entities.TypeCategory}

	err := s.repo.AddCategoryToDream(s.ctx, id, cat)

	s.NoError(err)

	var count int64
	err = s.db.Table("categories_dreams").Count(&count).Error
	s.NoError(err)
	s.EqualValues(1, count)
}

func (s *repoTestSuite) TestAddCategoryToDream_DreamNotFound() {
	cat := &entities.Category{Name: "Cat", Type: entities.TypeCategory}

	err := s.repo.AddCategoryToDream(s.ctx, 100, cat)

	s.ErrorIs(err, gorm.ErrForeignKeyViolated)
}

func (s *repoTestSuite) TestRemoveCategoryFromDream() {
	cat := &entities.Category{Name: "name", Type: entities.TypeCategory}
	s.db.Create(cat)
	dreamId := s.createDreamWith("", false, nil, []entities.Category{*cat})

	err := s.repo.RemoveCategoryFromDream(s.ctx, dreamId, cat.ID)
	s.NoError(err)

	count := s.db.
		Model(&entities.Dream{ID: dreamId}).
		Association("Categories").
		Count()
	s.EqualValues(0, count)
}

func (s *repoTestSuite) TestRemoveCategoryFromDream_DreamNotFound() {
	cat := &entities.Category{Name: "name", Type: entities.TypeCategory}
	s.db.Create(cat)
	s.createDreamWith("", false, nil, []entities.Category{*cat})

	err := s.repo.RemoveCategoryFromDream(s.ctx, 100, cat.ID)
	s.NoError(err)
}
