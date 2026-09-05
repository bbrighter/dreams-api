package statistics

import "context"

type StatisticsService struct {
	repo *StatisticsRepo
}

func NewStatisticsService(repo *StatisticsRepo) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) CountByMonth(ctx context.Context) ([]CountByCatAndMonth, []CountByMonth, error) {
	catByMonth, err := s.repo.CountByCategoryAndMonth(ctx)
	if err != nil {
		return []CountByCatAndMonth{}, []CountByMonth{}, nil
	}
	dreamByMonth, err := s.repo.CountByDreamAndMonth(ctx)
	return catByMonth, dreamByMonth, err
}
