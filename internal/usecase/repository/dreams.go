package repository

import (
	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DreamsRepo struct {
	db *gorm.DB
}

func NewDreamsRepo(db *gorm.DB) *DreamsRepo {
	return &DreamsRepo{db: db}
}

func (r *DreamsRepo) List(showAll bool, includes []entity.Includes) entity.Dreams {
	var dreams entity.Dreams
	tx := r.db.Model(&entity.Dream{})
	if len(includes) > 0 {
		tx.Preload("Categories")
	}
	if !showAll {
		tx.Where(&entity.Dream{Visible: true})
	}
	tx.Find(&dreams)
	return dreams
}

func (r *DreamsRepo) Get(id uint, showAll bool) (entity.Dream, error) {
	tx := r.db.Model(&entity.Dream{}).Preload(clause.Associations)
	if !showAll {
		tx.Where(&entity.Dream{Visible: true})
	}

	dream := entity.Dream{ID: id}
	if tx.First(&dream).RowsAffected == 0 {
		return dream, entity.ErrorNotFound
	}
	return dream, tx.Error
}

func (r *DreamsRepo) Create(dream entity.Dream) (uint, error) {
	err := r.db.Create(&dream).Error
	return dream.ID, err
}

func (r *DreamsRepo) Update(dream entity.Dream) error {
	tx := r.db.Model(&dream).Updates(&dream)
	if tx.RowsAffected == 0 {
		return entity.ErrorNotFound
	}
	return tx.Error
}

func (r *DreamsRepo) Delete(dream entity.Dream) (entity.Categories, error) {
	var cats = &entity.Categories{}
	if rowsAffected := r.db.Preload(clause.Associations).
		Find(&dream).RowsAffected; rowsAffected == 0 {
		return *cats, entity.ErrorNotFound
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		tx.Select(clause.Associations).Delete(&dream)
		c, err := removeCategoriesIfNeeded(tx, dream.Categories)
		cats = &c
		return err
	})
	return *cats, err
}

func (r *DreamsRepo) ToggleVisibility(dream entity.Dream) error {
	if rowsAffected := r.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return entity.ErrorNotFound
	}
	return r.db.Model(&dream).Update("Visible", !dream.Visible).Error
}

func (r *DreamsRepo) Finalize(dreamId uint) error {
	tx := r.db.Model(&entity.Dream{}).
		Where("id = ?", dreamId).
		UpdateColumn("Finalized", true)
	if tx.RowsAffected == 0 {
		return entity.ErrorNotFound
	}
	return tx.Error
}

func removeCategoriesIfNeeded(db *gorm.DB, cats entity.Categories) (entity.Categories, error) {
	var unusedCategories entity.Categories
	for _, category := range cats {
		var usedCategory entity.Category
		db.Where(&entity.Category{Name: category.Name}).Preload("Dreams").Find(&usedCategory)
		if len(usedCategory.Dreams) == 0 {
			unusedCategories = append(unusedCategories, category)
		}
	}
	var leftOverCategories entity.Categories
	if len(unusedCategories) > 0 {
		if err := db.Delete(&unusedCategories).Error; err != nil {
			return leftOverCategories, err
		}
	}
	db.Find(&leftOverCategories)
	return leftOverCategories, nil
}
