package categories

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entities"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	db *gorm.DB
}

func NewCategoriesRepo(db *gorm.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}

func (r *CategoriesRepo) ListCategories(ctx context.Context) ([]entities.Category, error) {
	return gorm.G[entities.Category](r.db).Find(ctx)
}

type CountByCategory struct {
	entities.Category
	Count int64 `gorm:"column:count"`
}

func (r *CategoriesRepo) ListAndCountCategories(ctx context.Context) ([]CountByCategory, error) {
	var counts []CountByCategory

	tx := r.db.WithContext(ctx).
		Table("categories_dreams").
		Joins("JOIN categories ON categories_dreams.category_id = categories.id").
		Select("id", "name", "type", "count(*) as count").
		Group("id, name, type")
	err := tx.Scan(&counts).Error
	return counts, err
}

func (r *CategoriesRepo) CreateCategory(ctx context.Context, cat *entities.Category) error {
	return gorm.G[entities.Category](r.db).Create(ctx, cat)
}

func (r *CategoriesRepo) DeleteCategory(ctx context.Context, id uint) error {
	rows, err := gorm.G[entities.Category](r.db).
		Where("id = ?", id).
		Delete(ctx)

	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CategoriesRepo) UpdateCategory(ctx context.Context, id uint, updates map[string]any) error {
	tx := r.db.Model(&entities.Category{}).
		Where("id = ?", id).
		Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CategoriesRepo) MergeCategories(
	ctx context.Context,
	sourceCategoryId uint,
	targetCategoryId uint,
	newName string,
) error {
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

		if err := tx.Model(entities.Categories{}).
			Where("id = ?", targetCategoryId).
			// Where("type = ?", entities.TypeCategory).
			Update("name", newName).
			Error; err != nil {
			return nil
		}

		return tx.
			Where("id = ?", sourceCategoryId).
			Delete(&entities.Categories{}).
			Error
	})
}
