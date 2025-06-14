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

	// // Test one dream and tag exists
	// CreateTestDream(1, 1, true, t)
	// CreateTestDream(1, 1, false, t)

	// dreams = repo.GetDreams(&showAll)

	// assert.Len(t, dreams, 2)
	// d := dreams[1]
	// assert.Len(t, d.Categories, 0)
	// assert.Len(t, d.Persons, 0)

	// showAll = true
	// dreams = repo.GetDreams(&showAll)
	// assert.Len(t, dreams, 3)
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

func TestFinalize(t *testing.T) {
	repo := setupDreamsTest(t)
	err := repo.db.Create(&entity.Dream{ID: 1}).Error
	assert.NoError(t, err)

	err = repo.Finalize(1)
	assert.NoError(t, err)

	err = repo.Finalize(100)
	assert.Error(t, err)
}

func TestCreateDream(t *testing.T) {
	repo := setupDreamsTest(t)
	id, err := repo.Create(entity.Dream{Date: time.Now()})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
}

func TestUpdateDream(t *testing.T) {
	repo := setupDreamsTest(t)
	err := repo.db.Create(&entity.Dream{ID: 1}).Error
	assert.NoError(t, err)

	err = repo.Update(entity.Dream{ID: 1, Date: time.Now(), Description: "desc"})
	assert.NoError(t, err)

	err = repo.Update(entity.Dream{ID: 100})
	assert.ErrorIs(t, err, entity.ErrorNotFound)
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
