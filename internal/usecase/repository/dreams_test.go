package repository

import (
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupDreamsTest(t *testing.T) *DreamsRepo {
	logger, _ := zap.NewDevelopment()
	db := NewDatabase(":memory:", logger)
	repo := NewDreamsRepo(db)
	err := repo.db.AutoMigrate(
		&entity.Dream{},
		&entity.Category{},
	)
	assert.NoError(t, err)

	return repo
}

func TestGetDreams(t *testing.T) {
	tests := map[string]struct {
		createDreamBefore          bool
		includePersons             bool
		includeCategories          bool
		expectedLength             int
		expectedNumberOfCategories int
		expectedNumberOfPersons    int
	}{
		"no dreams":              {},
		"one dream":              {createDreamBefore: true, expectedLength: 1},
		"one dream + cat":        {createDreamBefore: true, includeCategories: true, expectedLength: 1, expectedNumberOfCategories: 1, expectedNumberOfPersons: 1},
		"one dream + pers":       {createDreamBefore: true, includePersons: true, expectedLength: 1, expectedNumberOfPersons: 1, expectedNumberOfCategories: 1},
		"one dream + cat + pers": {createDreamBefore: true, includePersons: true, includeCategories: true, expectedLength: 1, expectedNumberOfCategories: 1, expectedNumberOfPersons: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo := setupDreamsTest(t)
			if test.createDreamBefore {
				repo.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{
					entity.Category{ID: 1, Type: entity.TypeCategory},
					entity.Category{ID: 2, Type: entity.TypePerson},
				}})
			}

			includes := []entity.Includes{}
			if test.includeCategories {
				includes = append(includes, entity.IncludeCategories)
			}
			if test.includePersons {
				includes = append(includes, entity.IncludePersons)
			}
			dreams := repo.List(false, includes)
			assert.Len(t, dreams, test.expectedLength)
			if test.expectedLength > 0 {
				numberOfPersons := 0
				numberOfCats := 0
				for _, cat := range dreams[0].Categories {
					if cat.Type == entity.TypeCategory {
						numberOfCats++
					}
					if cat.Type == entity.TypePerson {
						numberOfPersons++
					}
				}
				assert.EqualValues(t, numberOfCats, test.expectedNumberOfCategories)
				assert.EqualValues(t, numberOfPersons, test.expectedNumberOfPersons)
			}
		})
	}
}

func TestGetDream(t *testing.T) {
	repo := setupDreamsTest(t)

	var dream entity.Dream
	var err error

	// Test no dreams exist
	_, err = repo.Get(1, false)

	assert.Error(t, err)

	// Test dream exists
	repo.db.Create(&entity.Dream{ID: 100})
	dream, err = repo.Get(100, false)

	assert.NoError(t, err)
	assert.EqualValues(t, 100, dream.ID)
}

func TestCreateDream(t *testing.T) {
	repo := setupDreamsTest(t)
	id, err := repo.Create(entity.Dream{Date: time.Now()})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
}

func TestUpdateDream(t *testing.T) {
	tests := map[string]struct {
		dreamId   uint
		updateMap map[string]any
		isError   bool
	}{
		"date, ok":        {dreamId: 1, updateMap: map[string]any{"date": time.Now()}},
		"description, ok": {dreamId: 1, updateMap: map[string]any{"description": "desc"}},
		"date + desc, ok": {dreamId: 1, updateMap: map[string]any{"date": time.Now(), "description": "desc"}},
		"rating, ok":      {dreamId: 1, updateMap: map[string]any{"rating": 3}},
		"finalized, ok":   {dreamId: 1, updateMap: map[string]any{"finalized": true}},
		"not found":       {dreamId: 100, updateMap: map[string]any{"date": time.Now()}, isError: true},
		"invalid column":  {dreamId: 1, updateMap: map[string]any{"invalid": "abc"}, isError: true},
		"invalid type":    {dreamId: 1, updateMap: map[string]any{"rating": "a"}, isError: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var err error
			repo := setupDreamsTest(t)
			err = repo.db.Create(&entity.Dream{ID: 1}).Error
			assert.NoError(t, err)

			err = repo.Update(test.dreamId, test.updateMap)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})

	}

	// err = repo.Update(entity.Dream{ID: 1, Date: time.Now(), Description: "desc"})
	// assert.NoError(t, err)

	// err = repo.Update(entity.Dream{ID: 100})
	// assert.ErrorIs(t, err, entity.ErrorNotFound)
}

func TestDeleteDream(t *testing.T) {
	repo := setupDreamsTest(t)
	err := repo.db.Create(&entity.Dream{
		ID: 1,
		Categories: entity.Categories{
			entity.Category{ID: 1, Name: "Cat", Type: entity.TypeCategory},
			entity.Category{ID: 2, Name: "Person", Type: entity.TypePerson},
		},
	}).Error
	assert.NoError(t, err)

	cats, err := repo.Delete(entity.Dream{ID: 1})
	assert.NoError(t, err)
	assert.Len(t, cats, 0)

	var categories entity.Categories
	repo.db.Find(&categories)
	assert.Len(t, categories, 0)
}

func TestToggleVisibility(t *testing.T) {
	repo := setupDreamsTest(t)
	err := repo.db.Create(&entity.Dream{ID: 1}).Error
	assert.NoError(t, err)

	err = repo.ToggleVisibility(entity.Dream{ID: 1})
	assert.NoError(t, err)

	var dream entity.Dream
	repo.db.First(&dream)
	assert.False(t, dream.Visible)
}
