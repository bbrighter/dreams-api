package statistics

import "context"

type StatisticsService struct {
	repo *StatisticsRepo
}

func NewStatisticsService(repo *StatisticsRepo) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) CountByMonth(ctx context.Context) ([]CountByCatAndMonth, error) {
	return s.repo.CountByCategoryAndMonth(ctx)
}
