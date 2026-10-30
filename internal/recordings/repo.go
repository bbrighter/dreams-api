package recordings

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entities"
	"gorm.io/gorm"
)

type RecordingsRepository struct {
	db *gorm.DB
}

func NewRecordingsRepository(db *gorm.DB) *RecordingsRepository {
	return &RecordingsRepository{db: db}
}

func (r *RecordingsRepository) Create(ctx context.Context, rec *entities.Recording) error {
	return gorm.G[entities.Recording](r.db).Create(ctx, rec)
}

func (r *RecordingsRepository) Exists(ctx context.Context, dream_id uint) (bool, error) {
	count, err := gorm.G[entities.Recording](r.db).Where("dream_id = ?", dream_id).Count(ctx, "*")
	if err != nil {
		return true, err
	}
	return count > 0, nil
}

func (r *RecordingsRepository) First(ctx context.Context, id uint) (entities.Recording, error) {
	return gorm.G[entities.Recording](r.db).Where("id = ?", id).First(ctx)
}

func (r *RecordingsRepository) Find(ctx context.Context) ([]entities.Recording, error) {
	return gorm.G[entities.Recording](r.db).Find(ctx)
}

func (r *RecordingsRepository) Delete(ctx context.Context, id uint) error {
	rows, err := gorm.G[entities.Recording](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *RecordingsRepository) Update(ctx context.Context, id uint, values map[string]any) error {
	rows, err := gorm.G[map[string]any](r.db).Table("recordings").Where("id = ?", id).Updates(ctx, values)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
