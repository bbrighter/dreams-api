package repository

import (
	customerrors "github.com/bbrighter/dreams-api/internal/customErrors"
	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	Repo *gorm.DB
}

func NewCategoriesRepo(dbName string) *CategoriesRepo {
	db, _ := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	return &CategoriesRepo{Repo: db}
}

func (r *CategoriesRepo) GetAll() entity.Categories {
	var cats entity.Categories
	r.Repo.Find(&cats)
	return cats
}

func (r *CategoriesRepo) AddToDream(categoryName string, dream entity.Dream) (entity.Categories, error) {
	if rowsAffected := r.Repo.First(&dream).RowsAffected; rowsAffected == 0 {
		return nil, customerrors.ErrorNotFound
	}

	var category = entity.Category{
		Name:   categoryName,
		Dreams: entity.Dreams{dream},
	}
	r.Repo.Where(&entity.Category{Name: categoryName}).First(&category)
	err := r.Repo.Save(&category).Error

	var categories []entity.Category
	r.Repo.Find(&categories)
	return categories, err
}

func (r *CategoriesRepo) RemoveFromDream(category entity.Category, dream entity.Dream) (entity.Categories, error) {
	var categories = []entity.Category{}

	if rowsAffected := r.Repo.First(&dream).RowsAffected; rowsAffected == 0 {
		return categories, customerrors.ErrorNotFound
	}
	if rowsAffected := r.Repo.First(&category).RowsAffected; rowsAffected == 0 {
		return categories, customerrors.ErrorNotFound
	}

	r.Repo.Model(&dream).Association("Categories").Delete(category)
	categories, err := removeCategoriesIfNeeded(r.Repo, entity.Categories{category})

	return categories, err
}
