package statistics

import (
	"context"

	"gorm.io/gorm"
)

type StatisticsRepo struct {
	db *gorm.DB
}

func NewStatisticsRepo(db *gorm.DB) *StatisticsRepo {
	return &StatisticsRepo{db: db}
}

func (r StatisticsRepo) CountByCategoryAndMonth(ctx context.Context) ([]CountByCatAndMonth, error) {
	var aggs []CountByCatAndMonth

	err := r.db.WithContext(ctx).
		Model("dreams").
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

func (r *StatisticsRepo) CountByDreamAndMonth(ctx context.Context) ([]CountByMonth, error) {
	var aggs []CountByMonth

	err := r.db.WithContext(ctx).
		Model("dreams").
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
