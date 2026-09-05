package categories

import (
	"context"
	"testing"

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
	cats *CategoriesRepo
}

func (s *repoTestSuite) SetupSuite() {
	logger, _ := zap.NewDevelopment()
	s.db = migrations.NewDatabase(":memory:", logger)
	s.ctx = context.Background()
	s.cats = NewCategoriesRepo(s.db)
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

func (s *repoTestSuite) createCategory() uint {
	return s.createCategoryWith("cat", entities.TypeCategory)
}

func (s *repoTestSuite) createCategoryWith(name string, catType entities.CategoryType) uint {
	var cat = &entities.Category{Name: name, Type: catType}
	err := gorm.G[entities.Category](s.db).Create(s.ctx, cat)
	s.Require().NoError(err)
	return cat.ID
}

func (s *repoTestSuite) createDreamWithCat(catId uint) uint {
	var dream = &entities.Dream{Categories: entities.Categories{entities.Category{ID: catId}}}
	err := gorm.G[entities.Dream](s.db).Create(s.ctx, dream)
	s.Require().NoError(err)
	return dream.ID
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

func (s *repoTestSuite) TestListCategoriesNoEntries() {
	cats, err := s.cats.ListCategories(s.ctx)
	s.NoError(err)
	s.Len(cats, 0)
}

func (s *repoTestSuite) TestListCategories() {
	catId := s.createCategory()

	cats, err := s.cats.ListCategories(s.ctx)
	s.NoError(err)
	s.Len(cats, 1)
	s.Equal(cats[0].Name, "cat")
	s.Equal(cats[0].ID, catId)
	s.Equal(cats[0].Type, entities.TypeCategory)
}

func (s *repoTestSuite) TestListAndCountCategories() {
	catId := s.createCategory()
	s.createDreamWithCat(catId)

	cats, err := s.cats.ListAndCountCategories(s.ctx)
	s.NoError(err)
	s.Len(cats, 1)
	s.EqualValues(cats[0].Count, 1)
	s.Equal(cats[0].ID, catId)
	s.Equal(cats[0].Name, "cat")
	s.Equal(cats[0].Type, entities.TypeCategory)
}

func (s *repoTestSuite) TestCreateCategory() {
	cat := &entities.Category{Name: "name", Type: entities.TypePerson}
	err := s.cats.CreateCategory(s.ctx, cat)

	s.NoError(err)
	dbCat, err := gorm.G[entities.Category](s.db).First(s.ctx)
	s.Equal(*cat, dbCat)
}

func (s *repoTestSuite) TestDeleteCategoryNotFound() {
	err := s.cats.DeleteCategory(s.ctx, 100)

	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *repoTestSuite) TestDeleteCategory() {
	catId := s.createCategory()

	err := s.cats.DeleteCategory(s.ctx, catId)

	s.NoError(err)
	count, _ := gorm.G[entities.Category](s.db).Count(s.ctx, "*")
	s.EqualValues(0, count)
}

func (s *repoTestSuite) TestUpdateCategories() {
	catId := s.createCategory()

	var updates = map[string]any{
		"name": "new name",
		"type": entities.TypePerson,
	}
	err := s.cats.UpdateCategory(s.ctx, catId, updates)

	s.NoError(err)
	dbCat, _ := gorm.G[entities.Category](s.db).First(s.ctx)
	s.Equal("new name", dbCat.Name)
	s.Equal(entities.TypePerson, dbCat.Type)
}

func (s *repoTestSuite) TestUpdateCategoriesNotFound() {
	err := s.cats.UpdateCategory(s.ctx, 100, make(map[string]any))

	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *repoTestSuite) TestUpdateInvalidColumn() {
	catId := s.createCategory()

	var updates = map[string]any{
		"invalid": 123,
	}
	err := s.cats.UpdateCategory(s.ctx, catId, updates)
	s.Error(err)
}

func (s *repoTestSuite) TestMergeCategories() {
	sourceCatId := s.createCategory()
	targetCatId := s.createCategoryWith("person", entities.TypePerson)
	s.createDreamWithCat(sourceCatId)
	s.createDreamWithCat(targetCatId)

	err := s.cats.MergeCategories(s.ctx, sourceCatId, targetCatId, "new name")
	s.NoError(err)
	cats, _ := gorm.G[entities.Category](s.db).Find(s.ctx)
	s.Len(cats, 1)
	s.Equal("new name", cats[0].Name)
	s.Equal(entities.TypePerson, cats[0].Type)
	var count int64
	s.db.Table("categories_dreams").
		Where("category_id = ?", targetCatId).
		Count(&count)
	s.EqualValues(2, count)

	s.db.Table("categories_dreams").
		Where("category_id = ?", sourceCatId).
		Count(&count)
	s.EqualValues(0, count)
}
