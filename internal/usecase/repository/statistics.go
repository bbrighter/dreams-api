package repository

import (
	"sort"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type StatisticsRepo struct {
	Repo *gorm.DB
}

func NewStatisticsRepo(db *gorm.DB) *StatisticsRepo {
	return &StatisticsRepo{Repo: db}
}

func (r StatisticsRepo) CountCategories(showAll bool, maxNumberOfResults int) entity.Counts {
	tx := r.Repo.Preload("Categories")
	if !showAll {
		tx.Where(&entity.Dream{Visible: true})
	}
	var dreams entity.Dreams
	tx.Find(&dreams)

	var collection = make(map[uint]int)
	for _, dream := range dreams {
		for _, cat := range dream.Categories {
			_, exists := collection[cat.ID]
			if exists {
				collection[cat.ID] += 1
			} else {
				collection[cat.ID] = 1
			}
		}
	}

	var categoryCount entity.Counts
	for id, count := range collection {
		count := entity.Count{ID: id, Count: count}
		categoryCount = append(categoryCount, count)
	}
	sort.Slice(categoryCount, func(i, j int) bool {
		return categoryCount[i].Count > categoryCount[j].Count
	})

	if maxNumberOfResults < 1 {
		return categoryCount
	}
	var limit int = min(len(categoryCount), maxNumberOfResults)
	return categoryCount[:limit]
}

func (r StatisticsRepo) CountPersons(showAll bool, maxNumberOfResults int) entity.Counts {
	tx := r.Repo.Preload("Persons")
	if !showAll {
		tx.Where(&entity.Dream{Visible: true})
	}
	var dreams entity.Dreams
	tx.Find(&dreams)

	var categoryCount entity.Counts
	for _, dream := range dreams {
		for _, cat := range dream.Persons {
			index, exists := idInCategoryCount(cat.ID, categoryCount)
			if exists {
				categoryCount[index].Count += 1
			} else {
				count := entity.Count{ID: cat.ID, Count: 1}
				categoryCount = append(categoryCount, count)
			}
		}
	}
	sort.Slice(categoryCount, func(i, j int) bool {
		return categoryCount[i].Count > categoryCount[j].Count
	})

	if maxNumberOfResults < 1 {
		return categoryCount
	}
	var limit int = min(len(categoryCount), maxNumberOfResults)
	return categoryCount[:limit]
}

func idInCategoryCount(id uint, results entity.Counts) (int, bool) {
	for i, r := range results {
		if id == r.ID {
			return i, true
		}
	}
	return -1, false
}
