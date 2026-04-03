package repository

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type DreamsRepo struct {
	db *gorm.DB
}

func NewDreamsRepo(db *gorm.DB) *DreamsRepo {
	return &DreamsRepo{db: db}
}

func DreamIsNotHidden(tx *gorm.Statement) {
	tx.Where("visible = ?", true)
}

func DreamById(id uint) func(tx *gorm.Statement) {
	return func(tx *gorm.Statement) {
		tx.Where("id = ?", id)
	}
}

func (r *DreamsRepo) List(ctx context.Context, showAll bool) (entity.Dreams, error) {
	tx := gorm.G[entity.Dream](r.db).Where("1=1").Preload("Categories", nil)
	if !showAll {
		tx = tx.Scopes(DreamIsNotHidden)
	}
	return tx.Find(ctx)
}

func (r *DreamsRepo) Get(ctx context.Context, id uint, showAll bool) (entity.Dream, error) {
	tx := gorm.G[entity.Dream](r.db).
		Scopes(DreamById(id)).
		Preload("Categories", nil)
	if !showAll {
		tx = tx.Scopes(DreamIsNotHidden)
	}
	return tx.First(ctx)
}

func (r *DreamsRepo) Create(ctx context.Context, dream *entity.Dream) (uint, error) {
	err := gorm.G[entity.Dream](r.db).Create(ctx, dream)
	return dream.ID, err
}

func (r *DreamsRepo) Update(ctx context.Context, dreamId uint, updates map[string]any) error {
	if err := validateTypes(entity.Dream{}, updates); err != nil {
		return err
	}
	tx := r.db.Model(&entity.Dream{}).Where("id = ?", dreamId).Updates(updates)
	if tx.RowsAffected == 0 {
		return entity.ErrorNotFound
	}
	return tx.Error
}

func (r *DreamsRepo) Delete(ctx context.Context, id uint) error {
	tx := r.db.WithContext(ctx).
		Select("Categories").
		Delete(&entity.Dream{ID: id})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *DreamsRepo) ToggleVisibility(ctx context.Context, id uint) error {
	rows, err := gorm.G[entity.Dream](r.db).Scopes(DreamById(id)).Update(ctx, "Visible", gorm.Expr("NOT visible"))
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
