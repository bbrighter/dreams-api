package repository

import (
	customerrors "github.com/bbrighter/dreams-api/internal/customErrors"
	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DreamsRepo struct {
	Repo *gorm.DB
}

func NewDreamsRepo(dbName string) *DreamsRepo {
	db, _ := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	return &DreamsRepo{Repo: db}
}

func (r *DreamsRepo) GetAll(showAll bool) entity.Dreams {
	var dreams entity.Dreams
	tx := r.Repo.Model(&entity.Dream{})
	if !showAll {
		tx.Where(&entity.Dream{Visible: true})
	}
	tx.Find(&dreams)
	return dreams
}

func (r *DreamsRepo) GetById(id uint, showAll bool) (entity.Dream, error) {
	tx := r.Repo.Model(&entity.Dream{}).Preload(clause.Associations)
	if !showAll {
		tx.Where(&entity.Dream{Visible: true})
	}

	dream := entity.Dream{ID: id}
	if tx.First(&dream).RowsAffected == 0 {
		return dream, customerrors.ErrorNotFound
	}
	return dream, tx.Error
}

func (r *DreamsRepo) Create(dream entity.Dream) (uint, error) {
	if err := r.Repo.Create(&dream).Error; err != nil {
		return 0, err
	}
	return dream.ID, nil
}

func (r *DreamsRepo) Update(dream entity.Dream) error {
	tx := r.Repo.Model(&dream).Updates(&dream)
	if tx.RowsAffected == 0 {
		return customerrors.ErrorNotFound
	}
	return tx.Error
}

func (r *DreamsRepo) Delete(dream entity.Dream) (entity.Categories, error) {
	if rowsAffected := r.Repo.Preload(clause.Associations).
		Find(&dream).RowsAffected; rowsAffected == 0 {
		return entity.Categories{}, customerrors.ErrorNotFound
	}
	r.Repo.Select(clause.Associations).Delete(&dream)
	return removeCategoriesIfNeeded(r.Repo, dream.Categories)

}

func (r *DreamsRepo) ToggleVisibility(dream entity.Dream) error {
	if rowsAffected := r.Repo.First(&dream).RowsAffected; rowsAffected == 0 {
		return customerrors.ErrorNotFound
	}

	if err := r.Repo.Model(&dream).Update("Visible", !dream.Visible).Error; err != nil {
		return err
	}
	return nil
}

func removeCategoriesIfNeeded(db *gorm.DB, cats entity.Categories) (entity.Categories, error) {
	for _, category := range cats {
		var usedCategory entity.Category
		db.Where(&entity.Category{Name: category.Name}).Preload("Dreams").Find(&usedCategory)
		if len(usedCategory.Dreams) == 0 {
			db.Delete(&category)
		}
	}
	var leftOverCategories entity.Categories
	db.Find(&leftOverCategories)
	return leftOverCategories, nil
}
