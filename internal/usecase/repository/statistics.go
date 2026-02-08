package repository

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type StatisticsRepo struct {
	db *gorm.DB
}

func NewStatisticsRepo(db *gorm.DB) *StatisticsRepo {
	return &StatisticsRepo{db: db}
}

func (r *StatisticsRepo) CountByCategory(ctx context.Context, limit int) (entity.CategoriesCount, error) {
	var aggs entity.CategoriesCount

	tx := r.db.WithContext(ctx).
		Table("categories_dreams").
		Select("category_id", "count(*) as count").
		Group("category_id")
	if limit > 0 {
		tx = tx.Order("count desc").Limit(limit)
	}
	err := tx.Scan(&aggs).Error
	return aggs, err
}

func (r StatisticsRepo) CountByCategoryAndMonth(ctx context.Context) ([]entity.CountByCatAndMonth, error) {
	var aggs []entity.CountByCatAndMonth

	err := r.db.Debug().WithContext(ctx).
		Model(&entity.Dream{}).
		Select(
			"categories.id as category_id",
			"strftime('%Y/%m', dreams.date) as month",
			"count(dreams.id) as count",
		).
		Joins("LEFT JOIN categories_dreams on dreams.id = categories_dreams.dream_id").
		Joins("LEFT JOIN categories on categories_dreams.category_id = categories.id").
		Group("categories.id").Group("strftime('%Y/%m', dreams.date)").
		Scan(&aggs).Error
	return aggs, err
}

func (r *StatisticsRepo) CountByDreamAndMonth(ctx context.Context) ([]entity.CountByMonth, error) {
	var aggs []entity.CountByMonth

	err := r.db.WithContext(ctx).
		Model(&entity.Dream{}).
		Select(
			"strftime('%Y/%m', dreams.date) as month",
			"count(distinct dreams.id) as count",
		).
		Joins("LEFT JOIN categories_dreams on dreams.id = categories_dreams.dream_id").
		Joins("LEFT JOIN categories on categories_dreams.category_id = categories.id").
		Group("strftime('%Y/%m', dreams.date)").
		Scan(&aggs).Error
	return aggs, err
}
