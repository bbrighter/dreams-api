package repository

import (
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func setupCategoriesTest(t *testing.T) *CategoriesRepo {
	logger, _ := zap.NewDevelopment()
	db := NewDatabase(":memory:", logger)
	repo := NewCategoriesRepo(db)
	err := repo.db.AutoMigrate(
		&entity.Dream{},
		&entity.Category{},
	)
	assert.NoError(t, err)

	return repo
}

func (s *RepoTestSuite) TestGetAllCategories() {
	tests := map[string]struct {
		createCats    bool
		expectedCount int
	}{
		"0 found": {},
		"1 found": {createCats: true, expectedCount: 1},
	}
	for name, test := range tests {
		s.Run(name, func() {
			if test.createCats {
				s.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{entity.Category{ID: 1, Type: entity.TypeCategory}}})
			}

			cats, err := s.cats.List(s.ctx)
			s.NoError(err)
			s.Len(cats, test.expectedCount)
		})
	}
}

func (s *RepoTestSuite) TestRemoveCategoryFromDream() {
	tests := map[string]struct {
		dreamExists bool
		catExits    bool
		expectError error
	}{
		"no dream":    {catExits: true, expectError: nil},
		"no cat":      {dreamExists: true, expectError: nil},
		"dream + cat": {dreamExists: true, catExits: true},
	}

	for name, test := range tests {
		s.Run(name, func() {
			var cat = entity.Category{ID: 10, Name: "name", Type: entity.TypeCategory}
			var dream = entity.Dream{ID: 1}
			if test.catExits && !test.dreamExists {
				err := gorm.G[entity.Category](s.db).Create(s.ctx, &cat)
				s.Require().NoError(err)
			}
			if test.dreamExists && !test.catExits {
				err := gorm.G[entity.Dream](s.db).Create(s.ctx, &dream)
				s.Require().NoError(err)
			}
			if test.dreamExists && test.catExits {
				dream.Categories = entity.Categories{cat}
				err := gorm.G[entity.Dream](s.db).Create(s.ctx, &dream)
				s.Require().NoError(err)
			}

			err := s.cats.RemoveFromDream(s.ctx, 1, 10)
			if test.expectError != nil {
				s.Equal(test.expectError, err)
				return
			}
			s.NoError(err)
		})
	}
}

func (s *RepoTestSuite) TestCreateCategory() {
	tests := map[string]struct {
		categories    []entity.Category
		expectedError error
	}{
		"ok": {categories: []entity.Category{{Name: "name", Type: entity.TypeCategory}}},
		"same name, different type": {categories: []entity.Category{
			{Name: "name", Type: entity.TypeCategory},
			{Name: "name", Type: entity.TypePerson},
		}},
		"same name, same type": {categories: []entity.Category{
			{Name: "name", Type: entity.TypeCategory},
			{Name: "name", Type: entity.TypeCategory},
		}, expectedError: gorm.ErrDuplicatedKey},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var err error
			for _, cat := range test.categories {
				_, err = s.cats.Create(s.ctx, cat.Name, cat.Type)
			}
			if test.expectedError != nil {
				s.ErrorIs(err, test.expectedError)
				return
			}
			s.NoError(err)
		})
	}
}

func (s *RepoTestSuite) TestAddToDream() {
	tests := map[string]struct {
		dreamExists   bool
		catExists     bool
		expectedError error
	}{
		"cat and dream exist": {dreamExists: true, catExists: true},
		"not cat":             {dreamExists: true, expectedError: gorm.ErrRecordNotFound},
		"no dream":            {catExists: true, expectedError: gorm.ErrForeignKeyViolated},
	}
	for name, test := range tests {
		s.Run(name, func() {
			if test.dreamExists {
				err := gorm.G[entity.Dream](s.db).Create(s.ctx, &entity.Dream{ID: 1})
				s.Require().NoError(err)
			}
			if test.catExists {
				err := gorm.G[entity.Category](s.db).Create(s.ctx, &entity.Category{ID: 2, Type: entity.TypeCategory})
				s.Require().NoError(err)
			}

			err := s.cats.AddToDream(s.ctx, 1, 2)
			if test.expectedError != nil {
				s.ErrorIs(err, test.expectedError)
				return
			}
			s.NoError(err)

			count := s.db.Model(&entity.Dream{ID: 1}).Association("Categories").Count()
			s.EqualValues(count, 1)
		})
	}
}
