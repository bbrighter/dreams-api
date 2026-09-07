package statistics

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entities"
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

	monthly_counts := r.db.WithContext(ctx).
		Select(
			"'total' as result_type",
			"NULL as category_id",
			"strftime('%Y/%m', dreams.date) as month",
			"count(distinct dreams.id) as count",
		).Model(&entities.Dream{}).
		Group("month")

	category_monthly_counts := r.db.WithContext(ctx).Select(
		"'category' as result_type",
		"categories_dreams.category_id as category_id",
		"strftime('%Y/%m', dreams.date) as month",
		"count(dreams.id) as count",
	).Model(&entities.Dream{}).
		Joins("LEFT JOIN categories_dreams on dreams.id = categories_dreams.dream_id").
		Group("category_id").Group("month")

	err := r.db.
		Table("(? UNION ALL ?)", monthly_counts, category_monthly_counts).
		Scan(&aggs).Error

	return aggs, err
}
