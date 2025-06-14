package repository

import (
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func TestGetAllCategories(t *testing.T) {
	r := setupCategoriesTest(t)

	cats := r.List()
	assert.Len(t, cats, 0)

	r.db.Create(&entity.Category{ID: 1, Type: entity.TypeCategory})
	cats = r.List()
	assert.Len(t, cats, 1)
}

func TestAddCategoryToDream(t *testing.T) {
	r := setupCategoriesTest(t)
	var err error

	var dream = entity.Dream{ID: 1}
	_, err = r.AddToDream("name", dream, entity.TypeCategory)
	assert.Error(t, err)

	r.db.Create(&dream)

	var cats entity.Categories
	cats, err = r.AddToDream("name", dream, entity.TypeCategory)

	assert.NoError(t, err)
	assert.Len(t, cats, 1)

	cats, err = r.AddToDream("name", dream, entity.TypeCategory)

	assert.NoError(t, err)
	assert.Len(t, cats, 1)

	cats, err = r.AddToDream("new name", dream, entity.TypeCategory)

	assert.NoError(t, err)
	assert.Len(t, cats, 2)
}

func TestRemoveCategoryFromDream(t *testing.T) {
	tests := map[string]struct {
		exists               bool
		catInUseByOtherDream bool
		expectError          error
		expectedLengthOfCats int
	}{
		"no dream and no category": {expectError: entity.ErrorNotFound},
		"dream + cat":              {exists: true},
		"cat in use, not removed":  {catInUseByOtherDream: true, expectedLengthOfCats: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := setupCategoriesTest(t)
			var cat = entity.Categories{entity.Category{ID: 10, Name: "name", Type: entity.TypeCategory}}
			if test.exists || test.catInUseByOtherDream {
				var dream = entity.Dream{ID: 1, Categories: cat}
				r.db.Create(&dream)
			}
			if test.catInUseByOtherDream {
				var otherDream = entity.Dream{ID: 2, Categories: cat}
				r.db.Debug().Create(&otherDream)
			}
			cats, err := r.RemoveFromDream(entity.Category{ID: 10}, entity.Dream{ID: 1})
			assert.Equal(t, test.expectError, err)
			assert.Len(t, cats, test.expectedLengthOfCats)
		})
	}
}
