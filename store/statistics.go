package store

import "sort"

type Count struct {
	ID    uint
	Count int
}

func (repo Repo) CountCategories(showAll bool) []Count {
	tx := repo.db.Preload("Categories")
	if !showAll {
		visible := true
		tx.Where(&Dream{Visible: &visible})
	}
	var dreams []Dream
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

	var categoryCount []Count
	for id, count := range collection {
		count := Count{ID: id, Count: count}
		categoryCount = append(categoryCount, count)
	}
	sort.Slice(categoryCount, func(i, j int) bool {
		return categoryCount[i].Count > categoryCount[j].Count
	})

	var maxNumberOfResults = 20
	var limit int = min(len(categoryCount), maxNumberOfResults)

	return categoryCount[:limit]
}

func (repo Repo) CountPersons(showAll bool) []Count {
	tx := repo.db.Preload("Persons")
	if !showAll {
		visible := true
		tx.Where(&Dream{Visible: &visible})
	}
	var dreams []Dream
	tx.Find(&dreams)

	var categoryCount []Count
	for _, dream := range dreams {
		for _, cat := range dream.Persons {
			index, exists := idInCategoryCount(cat.ID, categoryCount)
			if exists {
				categoryCount[index].Count += 1
			} else {
				count := Count{ID: cat.ID, Count: 1}
				categoryCount = append(categoryCount, count)
			}
		}
	}
	sort.Slice(categoryCount, func(i, j int) bool {
		return categoryCount[i].Count > categoryCount[j].Count
	})

	var maxNumberOfResults = 20
	var limit int = min(len(categoryCount), maxNumberOfResults)
	return categoryCount[:limit]
}

func idInCategoryCount(id uint, results []Count) (int, bool) {
	for i, r := range results {
		if id == r.ID {
			return i, true
		}
	}
	return -1, false
}
