package main

import (
	"time"
)

type Dream struct {
	ID          uint
	Date        time.Time
	Description string
	Tags        []*Tag
}

func (repo Repo) getDreams() ([]Dream, error) {
	var dreams []Dream
	if err := repo.db.Find(&dreams).Error; err != nil {
		return nil, err
	}
	return dreams, nil
}

func (repo Repo) getDream(id uint) (Dream, error) {
	var dream Dream = Dream{ID: id}
	tx := repo.db.First(&dream)
	if tx.RowsAffected == 0 {
		return dream, ErrorNotFound
	}
	return dream, tx.Error
}

func (repo Repo) createDream(dream Dream) (uint, error) {
	if err := repo.db.Create(&dream).Error; err != nil {
		return 0, err
	}
	return dream.ID, nil
}

func (repo Repo) deleteDream(id uint) error {
	tx := repo.db.Delete(&Dream{ID: id})
	if tx.RowsAffected == 0 {
		return ErrorNotFound
	}
	return tx.Error
}
