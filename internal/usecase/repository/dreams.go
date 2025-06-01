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

func (r *DreamsRepo) List(showAll bool) entity.Dreams {
	var dreams entity.Dreams
	tx := r.db.Model(&entity.Dream{})
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
	if err := r.db.Create(&dream).Error; err != nil {
		return 0, err
	}
	return dream.ID, nil
}

func (r *DreamsRepo) Update(dream entity.Dream) error {
	tx := r.db.Model(&dream).Updates(&dream)
	if tx.RowsAffected == 0 {
		return entity.ErrorNotFound
	}
	return tx.Error
}

func (r *DreamsRepo) Delete(dream entity.Dream) (entity.Categories, entity.Persons, error) {
	var cats *entity.Categories
	var pers *entity.Persons
	if rowsAffected := r.db.Preload(clause.Associations).
		Find(&dream).RowsAffected; rowsAffected == 0 {
		return entity.Categories{}, entity.Persons{}, entity.ErrorNotFound
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		tx.Select(clause.Associations).Delete(&dream)
		c, err := removeCategoriesIfNeeded(tx, dream.Categories)
		cats = &c
		if err != nil {
			return err
		}
		p, err := removePersonsIfNeeded(tx, dream.Persons)
		pers = &p
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entity.Categories{}, entity.Persons{}, err
	}
	return *cats, *pers, nil
}

func (r *DreamsRepo) ToggleVisibility(dream entity.Dream) error {
	if rowsAffected := r.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return entity.ErrorNotFound
	}

	if err := r.db.Model(&dream).Update("Visible", !dream.Visible).Error; err != nil {
		return err
	}
	return nil
}

func (r *DreamsRepo) Finalize(dreamId uint) error {
	tx := r.db.Model(&entity.Dream{}).
		Where("id = ?", dreamId).
		UpdateColumn("finalized", true)
	if tx.RowsAffected == 0 {
		return entity.ErrorNotFound
	}
	return tx.Error
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

func removePersonsIfNeeded(db *gorm.DB, persons entity.Persons) (entity.Persons, error) {
	for _, p := range persons {
		var usedPerson entity.Person
		db.Where(&entity.Person{Name: p.Name}).Preload("Dreams").Find(&usedPerson)
		if len(usedPerson.Dreams) == 0 {
			db.Delete(&p)
		}
	}
	var leftOverPersons entity.Persons
	db.Find(&leftOverPersons)
	return leftOverPersons, nil
}
