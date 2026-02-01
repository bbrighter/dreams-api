package repository

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	db *gorm.DB
}

func NewCategoriesRepo(db *gorm.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}

func (r *CategoriesRepo) List(ctx context.Context) (entity.Categories, error) {
	return gorm.G[entity.Category](r.db).Find(ctx)
}

func (r *CategoriesRepo) WithTransaction(tx *gorm.DB) *CategoriesRepo {
	return NewCategoriesRepo(tx)
}

func (r *CategoriesRepo) Create(ctx context.Context, catName string, catType entity.CategoryType) (uint, error) {
	var cat = &entity.Category{Name: catName, Type: catType}
	err := gorm.G[entity.Category](r.db).Create(ctx, cat)
	return cat.ID, err
}

func (r *CategoriesRepo) AddToDream(ctx context.Context, dreamId uint, catId uint) error {
	cat, err := gorm.G[entity.Category](r.db).Where("id = ?", catId).First(ctx)
	if err != nil {
		return err
	}
	return r.db.Model(&entity.Dream{ID: dreamId}).
		Association("Categories").
		Append(&cat)
}

func (r *CategoriesRepo) RemoveFromDream(ctx context.Context, dreamID uint, catId uint) error {
	return r.db.Model(&entity.Dream{ID: dreamID}).
		Association("Categories").
		Delete(&entity.Category{ID: catId})
}

func (r *CategoriesRepo) FindByName(ctx context.Context, catName string, catType entity.CategoryType) (uint, error) {
	cat, err := gorm.G[entity.Category](r.db).Where("name = ?", catName).First(ctx)
	return cat.ID, err
}
