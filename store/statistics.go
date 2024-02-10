package store

type Count struct {
	ID    uint
	Count int
}

func (repo Repo) CountCategories() []Count {
	var categoryCount []Count
	repo.db.Table("categories_dreams").Select("count(*) as count, category_id as id").Group("category_id").Scan(&categoryCount)
	return categoryCount
}

func (repo Repo) CountPersons() []Count {
	var categoryCount []Count
	repo.db.Table("people_dreams").Select("count(*) as count, person_id as id").Group("person_id").Scan(&categoryCount)
	return categoryCount
}
