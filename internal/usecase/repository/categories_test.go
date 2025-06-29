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
	tests := map[string]struct {
		createCats          bool
		preload             bool
		expectedCount       int
		expectedCountDreams int
	}{
		"0 found":          {},
		"1 found":          {createCats: true, expectedCount: 1},
		"0 found, preload": {preload: true},
		"1 found, preload": {preload: true, createCats: true, expectedCount: 1, expectedCountDreams: 1},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := setupCategoriesTest(t)
			if test.createCats {
				r.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{entity.Category{ID: 1, Type: entity.TypeCategory}}})
			}

			var includes []entity.Includes
			if test.preload {
				includes = append(includes, entity.IncludeDreamsCount)
			}
			cats := r.List(includes)
			assert.Len(t, cats, test.expectedCount)
			if test.expectedCount > 0 {
				assert.Len(t, cats[0].Dreams, test.expectedCountDreams)
			}
		})
	}
}

func TestAddCategoryToDream(t *testing.T) {
	r := setupCategoriesTest(t)
	var err error

	var dream = entity.Dream{ID: 1}
	err = r.AddToDream("name", dream, entity.TypeCategory)
	assert.Error(t, err)

	r.db.Create(&dream)

	err = r.AddToDream("name", dream, entity.TypeCategory)
	assert.NoError(t, err)

	err = r.AddToDream("name", dream, entity.TypeCategory)

	assert.NoError(t, err)

	err = r.AddToDream("new name", dream, entity.TypeCategory)

	assert.NoError(t, err)
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
				r.db.Create(&otherDream)
			}
			err := r.RemoveFromDream(entity.Category{ID: 10}, entity.Dream{ID: 1})
			assert.Equal(t, test.expectError, err)
			var cats entity.Categories
			r.db.Find(&cats)
			assert.Len(t, cats, test.expectedLengthOfCats)
		})
	}
}

func TestCategoryUpdate(t *testing.T) {
	tests := map[string]struct {
		id           uint
		categoryType entity.CategoryType
		name         string
		expectedErr  error
	}{
		"not found": {
			id:           100,
			categoryType: entity.TypeCategory,
			expectedErr:  entity.ErrorNotFound,
		},
		"ok": {
			id:           1,
			categoryType: entity.TypeCategory,
			name:         "new name",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := setupCategoriesTest(t)
			r.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{entity.Category{ID: 1, Name: "cat", Type: entity.TypePerson}}})

			updates := make(map[string]any)
			if test.categoryType != "" {
				updates["type"] = test.categoryType
			}
			if test.name != "" {
				updates["name"] = test.name
			}
			err := r.Update(test.id, updates)
			if test.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err.Error(), test.expectedErr.Error())
			}
		})
	}
}

func TestCategoryDelete(t *testing.T) {
	tests := map[string]struct {
		catId       uint
		expectedErr error
	}{
		"ok":           {catId: 10},
		"not found":    {catId: 100, expectedErr: entity.ErrorNotFound},
		"still in use": {catId: 1, expectedErr: entity.ErrorBadParam},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := setupCategoriesTest(t)
			r.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{entity.Category{ID: 1, Name: "cat", Type: entity.TypePerson}}})
			r.db.Create(&entity.Category{ID: 10, Type: entity.TypeCategory})

			err := r.Delete(test.catId)
			if test.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.expectedErr.Error())
			}
		})
	}
}

func TestCountByNameAndType(t *testing.T) {
	tests := map[string]struct {
		name          string
		catType       entity.CategoryType
		expectedCount int64
	}{
		"1":                 {name: "cat", catType: entity.TypeCategory, expectedCount: 1},
		"0: different type": {name: "cat", catType: entity.TypePerson, expectedCount: 0},
		"0: different name": {name: "cat2", catType: entity.TypeCategory, expectedCount: 0},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := setupCategoriesTest(t)
			r.db.Create(&entity.Categories{
				entity.Category{ID: 1, Name: "cat", Type: entity.TypeCategory},
				entity.Category{ID: 3, Name: "pers", Type: entity.TypePerson},
			})

			count := r.CountByNameAndType(test.name, test.catType)
			assert.Equal(t, test.expectedCount, count)
		})
	}
}

func TestMerge(t *testing.T) {
	type CategoryDream struct {
		CategoryID uint
		DreamID    uint
	}

	var cat1 = entity.Category{ID: 1, Name: "cat1", Type: entity.TypeCategory}
	var cat2 = entity.Category{ID: 2, Name: "cat2", Type: entity.TypeCategory}
	var cat3 = entity.Category{ID: 3, Name: "cat3", Type: entity.TypeCategory}

	tests := map[string]struct {
		fromId             uint
		toId               uint
		newName            string
		expectedError      error
		expectedCategories entity.Categories
		expectedRelations  []CategoryDream
	}{
		"ok": {fromId: 1, toId: 3, newName: "new",
			expectedCategories: entity.Categories{{ID: 3, Name: "new", Type: entity.TypeCategory}, cat2},
			expectedRelations:  []CategoryDream{{CategoryID: 3, DreamID: 2}, {CategoryID: 2, DreamID: 1}, {CategoryID: 3, DreamID: 1}},
		},
		"ok, but both cats in same dream": {fromId: 1, toId: 2, newName: "new",
			expectedCategories: entity.Categories{{ID: 2, Name: "new", Type: entity.TypeCategory}, cat3},
			expectedRelations:  []CategoryDream{{CategoryID: 3, DreamID: 2}, {CategoryID: 2, DreamID: 1}},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := setupCategoriesTest(t)
			r.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{cat1, cat2}})
			r.db.Create(&entity.Dream{ID: 2, Categories: entity.Categories{cat3}})

			err := r.Merge(test.fromId, test.toId, test.newName)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.expectedError.Error())
			}
			var cats entity.Categories
			r.db.Find(&cats)
			assert.ElementsMatch(t, test.expectedCategories, cats)

			var actual []CategoryDream
			r.db.Table("categories_dreams").Find(&actual)
			assert.ElementsMatch(t, actual, test.expectedRelations)
		})

	}
}
