package usecase

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type Statistics interface {
	CountByMonth(ctx context.Context, showAll bool) ([]entity.CountByCatAndMonth, entity.CountByMonths, error)
	CountCategories(ctx context.Context, limit int) (entity.CategoriesCount, error)
}

type StatisticsUseCase struct {
	repo IStatisticsRepo
}

func NewStatisticsUseCase(r IStatisticsRepo) *StatisticsUseCase {
	return &StatisticsUseCase{
		repo: r,
	}
}

func (u StatisticsUseCase) CountByMonth(ctx context.Context, showAll bool) ([]entity.CountByCatAndMonth, entity.CountByMonths, error) {
	catByMonth, err := u.repo.CountByCategoryAndMonth(ctx)
	if err != nil {
		return []entity.CountByCatAndMonth{}, []entity.CountByMonth{}, err
	}
	dreamByMonth, err := u.repo.CountByDreamAndMonth(ctx)
	return catByMonth, dreamByMonth, err
}

func (u StatisticsUseCase) CountCategories(ctx context.Context, limit int) (entity.CategoriesCount, error) {
	return u.repo.CountByCategory(ctx, limit)
}
