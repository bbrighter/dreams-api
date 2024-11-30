package repository

import (
	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	repo *gorm.DB
}

func NewCategoriesRepo(db *gorm.DB) *CategoriesRepo {
	return &CategoriesRepo{repo: db}
}

func (r *CategoriesRepo) List() entity.Categories {
	var cats entity.Categories
	r.repo.Find(&cats)
	return cats
}

func (r *CategoriesRepo) AddToDream(categoryName string, dream entity.Dream) (entity.Categories, error) {
	if rowsAffected := r.repo.First(&dream).RowsAffected; rowsAffected == 0 {
		return nil, entity.ErrorNotFound
	}

	var category = entity.Category{
		Name:   categoryName,
		Dreams: entity.Dreams{dream},
	}
	r.repo.Where(&entity.Category{Name: categoryName}).First(&category)
	err := r.repo.Save(&category).Error

	var categories []entity.Category
	r.repo.Find(&categories)
	return categories, err
}

func (r *CategoriesRepo) RemoveFromDream(category entity.Category, dream entity.Dream) (entity.Categories, error) {
	var categories = []entity.Category{}

	if rowsAffected := r.repo.First(&dream).RowsAffected; rowsAffected == 0 {
		return categories, entity.ErrorNotFound
	}
	if rowsAffected := r.repo.First(&category).RowsAffected; rowsAffected == 0 {
		return categories, entity.ErrorNotFound
	}

	r.repo.Model(&dream).Association("Categories").Delete(category)
	categories, err := removeCategoriesIfNeeded(r.repo, entity.Categories{category})

	return categories, err
}
