package repository

import (
	customerrors "github.com/bbrighter/dreams-api/internal/customErrors"
	"github.com/bbrighter/dreams-api/internal/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	repo   *gorm.DB
	logger *zap.Logger
}

func NewCategoriesRepo(name string, logger *zap.Logger) *CategoriesRepo {
	db := newDatabase(name, logger)
	return &CategoriesRepo{repo: db, logger: logger}
}

func (r *CategoriesRepo) GetAll() entity.Categories {
	var cats entity.Categories
	r.repo.Find(&cats)
	return cats
}

func (r *CategoriesRepo) AddToDream(categoryName string, dream entity.Dream) (entity.Categories, error) {
	if rowsAffected := r.repo.First(&dream).RowsAffected; rowsAffected == 0 {
		return nil, customerrors.ErrorNotFound
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
		return categories, customerrors.ErrorNotFound
	}
	if rowsAffected := r.repo.First(&category).RowsAffected; rowsAffected == 0 {
		return categories, customerrors.ErrorNotFound
	}

	r.repo.Model(&dream).Association("Categories").Delete(category)
	categories, err := removeCategoriesIfNeeded(r.repo, entity.Categories{category})

	return categories, err
}
