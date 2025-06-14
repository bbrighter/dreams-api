package usecase

import "github.com/bbrighter/dreams-api/internal/entity"

type StatisticsUseCase struct {
	repo IStatisticsRepo
}

func NewStatisticsUseCase(r IStatisticsRepo) *StatisticsUseCase {
	return &StatisticsUseCase{
		repo: r,
	}
}

func (u StatisticsUseCase) GetStatistics(showAll bool, maxNumber int) (entity.Counts, entity.Counts) {
	return u.repo.CountCategories(showAll, maxNumber)
}
