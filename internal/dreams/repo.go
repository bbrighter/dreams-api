package dreams

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entities"
	"gorm.io/gorm"
)

type DreamsRepo struct {
	db *gorm.DB
}

func NewDreamsRepo(db *gorm.DB) *DreamsRepo {
	return &DreamsRepo{db: db}
}

func (r *DreamsRepo) ListDreams(ctx context.Context) (entities.Dreams, error) {
	tx := gorm.G[entities.Dream](r.db).Preload("Categories", nil)
	return tx.Find(ctx)
}

func (r *DreamsRepo) GetDream(ctx context.Context, id uint) (entities.Dream, error) {
	return gorm.G[entities.Dream](r.db).
		Where("id = ?", id).
		Preload("Categories", nil).
		First(ctx)
}

func (r *DreamsRepo) CreateDream(ctx context.Context, dream *entities.Dream) (uint, error) {
	err := gorm.G[entities.Dream](r.db).Create(ctx, dream)
	return dream.ID, err
}

func (r *DreamsRepo) UpdateDream(ctx context.Context, dreamId uint, updates map[string]any) error {
	tx := r.db.Model(&entities.Dream{}).Where("id = ?", dreamId).Updates(updates)
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return tx.Error
}

func (r *DreamsRepo) DeleteDream(ctx context.Context, id uint) error {
	rows, err := gorm.G[entities.Dream](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *DreamsRepo) AddCategoryToDream(ctx context.Context, dreamId uint, category *entities.Category) error {
	return r.db.Model(&entities.Dream{ID: dreamId}).
		Association("Categories").
		Append(category)
}

func (r *DreamsRepo) RemoveCategoryFromDream(ctx context.Context, dreamID uint, catId uint) error {
	return r.db.WithContext(ctx).
		Model(&entities.Dream{ID: dreamID}).
		Association("Categories").
		Delete(&entities.Category{ID: catId})
}
