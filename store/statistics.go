package store

type CategoryCount struct {
	CategoryID uint
	Count      int
}

func (repo Repo) CountCategories() []CategoryCount {
	var categoryCount []CategoryCount
	// repo.db.Debug().Model(&Category{}).Select("count(*) as count, id").Group("id").Scan(&categoryCount)
	repo.db.Debug().Table("categories_dreams").Select("count(*) as count, category_id").Group("category_id").Scan(&categoryCount)
	return categoryCount
}
