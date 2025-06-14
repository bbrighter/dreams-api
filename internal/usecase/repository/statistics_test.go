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
	err := repo.db.AutoMigrate(
		&entity.Dream{},
		&entity.Category{},
	)
	assert.NoError(t, err)

	return repo
}

// func TestCountPersons(t *testing.T) {
// 	r := setupStatisticsTest(t)

// 	counts := r.CountPersons(true, 100)
// 	assert.Len(t, counts, 0)

// 	r.db.Create(&entity.Dream{
// 		ID:         1,
// 		Date:       time.Now(),
// 		Visible:    true,
// 		Categories: entity.Categories{entity.Category{ID: 10, Type: entity.TypeCategory}, entity.Category{ID: 100, Type: entity.TypePerson}},
// 	})

// 	counts = r.CountPersons(true, 100)
// 	assert.Len(t, counts, 1)
// 	assert.Equal(t, counts[0].Count, 1)
// }

func TestCountCategories(t *testing.T) {
	r := setupStatisticsTest(t)

	catCounts, persCounts := r.CountCategories(true, 100)
	assert.Len(t, catCounts, 0)
	assert.Len(t, persCounts, 0)

	r.db.Create(&entity.Dream{
		ID:      1,
		Date:    time.Now(),
		Visible: true,
		Categories: entity.Categories{
			entity.Category{ID: 10, Type: entity.TypeCategory},
			entity.Category{ID: 11, Type: entity.TypeCategory},
			entity.Category{ID: 100, Type: entity.TypePerson},
		},
	})
	r.db.Create(&entity.Dream{
		ID:         2,
		Visible:    true,
		Categories: entity.Categories{entity.Category{ID: 10, Type: entity.TypeCategory}},
	})

	catCounts, persCounts = r.CountCategories(true, 100)
	assert.Len(t, catCounts, 2)
	assert.Len(t, persCounts, 1)

	for _, count := range catCounts {
		if count.ID == 10 {
			assert.Equal(t, count.Count, 2)
		} else if count.ID == 11 {
			assert.Equal(t, count.Count, 1)
		}
	}
}
