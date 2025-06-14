package repository

import (
	"sort"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type StatisticsRepo struct {
	db *gorm.DB
}

func NewStatisticsRepo(db *gorm.DB) *StatisticsRepo {
	return &StatisticsRepo{db: db}
}

func (r StatisticsRepo) CountCategories(showAll bool, maxNumberOfResults int) (entity.Counts, entity.Counts) {
	tx := r.db.Preload("Categories")
	if !showAll {
		tx.Where(&entity.Dream{Visible: true})
	}
	var dreams entity.Dreams
	tx.Find(&dreams)

	catCount := countCats(dreams, entity.TypeCategory)
	persCount := countCats(dreams, entity.TypePerson)

	if maxNumberOfResults < 1 {
		return entity.Counts{}, entity.Counts{}
	}

	var catLimit int = min(len(catCount), maxNumberOfResults)
	var persLimit int = min(len(persCount), maxNumberOfResults)
	return catCount[:catLimit], persCount[:persLimit]
}

func countCats(dreams entity.Dreams, catType entity.CategoryType) entity.Counts {
	collection := make(map[uint]int)
	for _, dream := range dreams {
		for _, cat := range dream.Categories {
			if cat.Type == catType {
				_, exists := collection[cat.ID]
				if exists {
					collection[cat.ID] += 1
				} else {
					collection[cat.ID] = 1
				}
			}
		}
	}
	var counts entity.Counts
	for id, count := range collection {
		count := entity.Count{ID: id, Count: count}
		counts = append(counts, count)
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})
	return counts
}
