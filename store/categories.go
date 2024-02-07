package store

type Category struct {
	ID     uint
	Name   string
	Dreams []Dream `gorm:"many2many:categories_dreams;"`
}

func (repo Repo) GetCategories() []Category {
	var categories []Category
	repo.db.Find(&categories)
	return categories
}

func (repo Repo) AddCategoryToDream(categoryName string, dreamId uint) ([]Category, error) {
	var dream Dream = Dream{ID: dreamId}
	if rowsAffected := repo.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return nil, ErrorNotFound
	}

	var category Category = Category{Name: categoryName, Dreams: []Dream{dream}}
	repo.db.Where(&Category{Name: categoryName}).First(&category)
	err := repo.db.Save(&category).Error

	var categories []Category
	repo.db.Find(&categories)
	return categories, err
}

func (repo Repo) RemoveCategoryFromDream(categoryId uint, dreamId uint) ([]Category, error) {
	var categories []Category = []Category{}

	var dream Dream = Dream{ID: dreamId}
	if rowsAffected := repo.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return categories, ErrorNotFound
	}
	var category Category = Category{ID: categoryId}
	if rowsAffected := repo.db.First(&category).RowsAffected; rowsAffected == 0 {
		return categories, ErrorNotFound
	}

	repo.db.Model(&dream).Association("Categories").Delete(category)
	categories = repo.removeCategoriesIfNeeded([]Category{category})

	return categories, nil
}

func (repo Repo) removeCategoriesIfNeeded(categories []Category) []Category {
	for _, category := range categories {
		var usedCategory Category
		repo.db.Where(&Category{Name: category.Name}).Preload("Dreams").Find(&usedCategory)
		if len(usedCategory.Dreams) == 0 {
			repo.db.Delete(category)
		}
	}
	repo.db.Find(&categories)
	return categories
}
