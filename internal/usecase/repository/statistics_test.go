package repository

import (
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupStatisticsTest(t *testing.T) *StatisticsRepo {
	logger, _ := zap.NewDevelopment()
	db := NewDatabase(":memory:", logger)
	repo := NewStatisticsRepo(db)
	err := repo.Repo.AutoMigrate(
		&entity.Dream{},
		&entity.Category{},
		&entity.Person{},
	)
	assert.NoError(t, err)

	return repo
}

func TestCountPersons(t *testing.T) {
	r := setupStatisticsTest(t)

	counts := r.CountPersons(true, 100)
	assert.Len(t, counts, 0)

	r.Repo.Create(&entity.Dream{
		ID:         1,
		Date:       time.Now(),
		Visible:    true,
		Categories: entity.Categories{entity.Category{ID: 10}},
		Persons:    entity.Persons{entity.Person{ID: 100}},
	})

	counts = r.CountPersons(true, 100)
	assert.Len(t, counts, 1)
	assert.Equal(t, counts[0].Count, 1)
}

func TestCountCategories(t *testing.T) {
	r := setupStatisticsTest(t)

	counts := r.CountCategories(true, 100)
	assert.Len(t, counts, 0)

	r.Repo.Create(&entity.Dream{
		ID:         1,
		Date:       time.Now(),
		Visible:    true,
		Categories: entity.Categories{entity.Category{ID: 10}, entity.Category{ID: 11}},
		Persons:    entity.Persons{entity.Person{ID: 100}},
	})
	r.Repo.Create(&entity.Dream{
		ID:         2,
		Visible:    true,
		Categories: entity.Categories{entity.Category{ID: 10}},
	})

	counts = r.CountCategories(true, 100)
	assert.Len(t, counts, 2)

	for _, count := range counts {
		if count.ID == 10 {
			assert.Equal(t, count.Count, 2)
		} else if count.ID == 11 {
			assert.Equal(t, count.Count, 1)
		}
	}
}
