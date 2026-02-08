package repository

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type ManagementRepo struct {
	db *gorm.DB
}

func NewManagementRepo(db *gorm.DB) *ManagementRepo {
	return &ManagementRepo{db: db}
}

func (r *ManagementRepo) List(ctx context.Context) (entity.CategoriesCount, error) {
	var cats entity.CategoriesCount
	err := r.db.WithContext(ctx).
		Select("categories.id", "count(*) as count").
		Table("categories_dreams").
		Joins("JOIN categories ON categories_dreams.category_id = categories.id").
		Group("categories.id").
		Scan(&cats).Error
	return cats, err
}

func (r *ManagementRepo) Update(ctx context.Context, categoryId uint, updates map[string]any) error {
	if err := validateTypes(entity.Category{}, updates); err != nil {
		return err
	}
	var cat = entity.Category{ID: categoryId}
	if rows := r.db.WithContext(ctx).First(&cat).RowsAffected; rows == 0 {
		return entity.ErrorNotFound
	}
	return r.db.WithContext(ctx).Model(&cat).Updates(updates).Error
}

func (r *ManagementRepo) Delete(ctx context.Context, categoryId uint) error {
	rows, err := gorm.G[entity.Category](r.db).Where("id = ?", categoryId).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ManagementRepo) Merge(ctx context.Context, sourceCategoryId, targetCategoryId uint, newName string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fromDreamIds []uint
		if err := tx.Table("categories_dreams").
			Where("category_id = ?", sourceCategoryId).
			Pluck("dream_id", &fromDreamIds).
			Error; err != nil {
			return err
		}

		if err := tx.Table("categories_dreams").
			Where("category_id = ?", targetCategoryId).
			Where("dream_id in ?", fromDreamIds).
			Delete(nil).
			Error; err != nil {
			return err
		}

		if err := tx.Table("categories_dreams").
			Where("category_id = ?", sourceCategoryId).
			UpdateColumn("category_id", targetCategoryId).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.Category{ID: targetCategoryId, Type: entity.TypeCategory}).
			Update("name", newName).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Category{ID: sourceCategoryId}).Error
	})
}
