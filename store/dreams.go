package store

import (
	"time"

	"gorm.io/gorm/clause"
)

type Dream struct {
	ID          uint
	Date        time.Time
	Description string
	Tags        []Tag    `gorm:"many2many:tags_dreams;"`
	Persons     []Person `gorm:"many2many:people_dreams;"`
}

func (repo Repo) GetDreams() []Dream {
	var dreams []Dream
	repo.db.Model(&Dream{}).Preload("Tags").Find(&dreams)
	return dreams
}

func (repo Repo) GetDream(id uint) (Dream, error) {
	var dream Dream = Dream{ID: id}
	tx := repo.db.Model(&Dream{}).Preload("Tags").First(&dream)
	if tx.RowsAffected == 0 {
		return dream, ErrorNotFound
	}
	return dream, tx.Error
}

func (repo Repo) CreateDream(dream Dream) (uint, error) {
	if err := repo.db.Create(&dream).Error; err != nil {
		return 0, err
	}
	return dream.ID, nil
}

func (repo Repo) UpdateDream(dream Dream) error {
	tx := repo.db.Model(&dream).Updates(&dream)
	if tx.RowsAffected == 0 {
		return ErrorNotFound
	}
	return tx.Error
}

func (repo Repo) DeleteDream(id uint) error {
	var dream Dream = Dream{ID: id}
	if rowsAffected := repo.db.Preload(clause.Associations).
		Find(&dream).RowsAffected; rowsAffected == 0 {
		return ErrorNotFound
	}
	repo.db.Select(clause.Associations).Delete(&dream)
	repo.removeTagsIfNeeded(dream.Tags)

	return nil
}
