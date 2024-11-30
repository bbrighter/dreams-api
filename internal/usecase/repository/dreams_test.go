package repository

import (
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupDreamsTest(t *testing.T) *DreamsRepo {
	logger, _ := zap.NewDevelopment()
	db := NewDatabase(":memory:", logger)
	repo := NewDreamsRepo(db)
	err := repo.Repo.AutoMigrate(
		&entity.Dream{},
		&entity.Category{},
		&entity.Person{},
	)
	assert.NoError(t, err)

	return repo
}

func TestGetDreams(t *testing.T) {
	repo := setupDreamsTest(t)

	var dreams []entity.Dream

	// Test no dreams exist
	dreams = repo.List(false)

	assert.Len(t, dreams, 0)

	// Test one dream exists
	repo.Repo.Create(&entity.Dream{ID: 1})
	dreams = repo.List(false)

	assert.Len(t, dreams, 1)

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
	repo.Repo.Create(&entity.Dream{ID: 100})
	dream, err = repo.Get(100, false)

	assert.NoError(t, err)
	assert.EqualValues(t, 100, dream.ID)
}
