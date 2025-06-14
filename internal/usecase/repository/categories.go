package repository

import (
	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	db *gorm.DB
}

func NewCategoriesRepo(db *gorm.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}

func (r *CategoriesRepo) List() entity.Categories {
	var cats entity.Categories
	r.db.Find(&cats)
	return cats
}

func (r *CategoriesRepo) AddToDream(categoryName string, dream entity.Dream, categoryType entity.CategoryType) (entity.Categories, error) {
	if rowsAffected := r.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return nil, entity.ErrorNotFound
	}

	var category = entity.Category{
		Name:   categoryName,
		Type:   categoryType,
		Dreams: entity.Dreams{dream},
	}
	r.db.Where(&entity.Category{Name: categoryName, Type: categoryType}).First(&category)
	err := r.db.Save(&category).Error

	var categories []entity.Category
	r.db.Find(&categories)
	return categories, err
}

func (r *CategoriesRepo) RemoveFromDream(category entity.Category, dream entity.Dream) (entity.Categories, error) {
	var categories = []entity.Category{}

	if rowsAffected := r.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return categories, entity.ErrorNotFound
	}
	if rowsAffected := r.db.First(&category).RowsAffected; rowsAffected == 0 {
		return categories, entity.ErrorNotFound
	}

	r.db.Model(&dream).Association("Categories").Delete(&category)
	categories, err := removeCategoriesIfNeeded(r.db, entity.Categories{category})

	return categories, err
}
