package repository

import (
	"strings"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type CategoriesRepo struct {
	db *gorm.DB
}

func NewCategoriesRepo(db *gorm.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}

func (r *CategoriesRepo) List(includes []entity.Includes) entity.Categories {
	var cats entity.Categories
	tx := r.db.Debug()
	if len(includes) > 0 {
		tx = tx.Preload("Dreams")
	}
	tx.Find(&cats)
	return cats
}

func (r *CategoriesRepo) AddToDream(categoryName string, dream entity.Dream, categoryType entity.CategoryType) error {
	if rowsAffected := r.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return entity.ErrorNotFound
	}

	categoryName = strings.TrimSpace(categoryName)
	var category = entity.Category{
		Name:   categoryName,
		Type:   categoryType,
		Dreams: entity.Dreams{dream},
	}
	r.db.Where(&entity.Category{Name: categoryName, Type: categoryType}).First(&category)
	return r.db.Save(&category).Error
}

func (r *CategoriesRepo) RemoveFromDream(category entity.Category, dream entity.Dream) error {
	if rowsAffected := r.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return entity.ErrorNotFound
	}
	if rowsAffected := r.db.First(&category).RowsAffected; rowsAffected == 0 {
		return entity.ErrorNotFound
	}

	r.db.Model(&dream).Association("Categories").Delete(&category)
	_, err := removeCategoriesIfNeeded(r.db, entity.Categories{category})

	return err
}

func (r *CategoriesRepo) Update(categoryId uint, updates map[string]any) error {
	if err := validateTypes(entity.Category{}, updates); err != nil {
		return err
	}
	var cat = entity.Category{ID: categoryId}
	if rows := r.db.First(&cat).RowsAffected; rows == 0 {
		return entity.ErrorNotFound
	}
	return r.db.Model(&cat).Updates(updates).Error
}

func (r *CategoriesRepo) CountByNameAndType(name string, categoryType entity.CategoryType) int64 {
	var numSameNames int64
	r.db.Model(&entity.Category{}).
		Where("name = ?", name).
		Where("type = ?", categoryType).
		Count(&numSameNames)
	return numSameNames
}

func (r *CategoriesRepo) Delete(categoryId uint) error {
	var count int64
	tx := r.db.Table("categories_dreams").
		Where("category_id = ?", categoryId).
		Count(&count)
	if tx.Error != nil {
		return tx.Error
	}
	if count > 0 {
		return entity.ErrorBadParamWithReasons("category still in use")
	}
	delTx := r.db.Delete(&entity.Category{ID: categoryId})
	if delTx.Error != nil {
		return delTx.Error
	}
	if delTx.RowsAffected == 0 {
		return entity.ErrorNotFound
	}
	return nil
}

func (r *CategoriesRepo) First(categoryId uint) (entity.Category, error) {
	var cat = entity.Category{ID: categoryId}
	if rows := r.db.First(&cat).RowsAffected; rows == 0 {
		return cat, entity.ErrorNotFound
	}
	return cat, nil
}

func (r *CategoriesRepo) Merge(sourceCategoryId, targetCategoryId uint, newName string) error {
	return r.db.Debug().Transaction(func(tx *gorm.DB) error {
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
