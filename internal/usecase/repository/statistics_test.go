package repository

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

func (s *RepoTestSuite) TestCountByCategory() {
	var cat1 = entity.Category{ID: 1, Type: entity.TypeCategory, Name: "Cat1"}
	var cat2 = entity.Category{ID: 2, Type: entity.TypeCategory, Name: "Cat2"}
	tests := map[string]struct {
		dreams            []entity.Dream
		expectedCat1Count int64
		expectedCat2Count int64
	}{
		"one cat, one dream": {
			expectedCat1Count: 1,
			dreams:            entity.Dreams{{Categories: entity.Categories{cat1}}},
		},
		"two cats, one dream": {
			expectedCat1Count: 1, expectedCat2Count: 1,
			dreams: entity.Dreams{{Categories: entity.Categories{cat1, cat2}}},
		},
		"one cat, two dreams": {
			expectedCat1Count: 2,
			dreams: entity.Dreams{
				entity.Dream{Categories: entity.Categories{cat1}},
				entity.Dream{Categories: entity.Categories{cat1}},
			},
		},
		"two cats, two dreams": {
			expectedCat1Count: 1, expectedCat2Count: 1,
			dreams: entity.Dreams{
				entity.Dream{Categories: entity.Categories{cat1}},
				entity.Dream{Categories: entity.Categories{cat2}},
			},
		},
	}
	for name, test := range tests {
		s.Run(name, func() {
			err := gorm.G[entity.Dream](s.db).CreateInBatches(s.ctx, &test.dreams, 100)
			s.Require().NoError(err)

			count, err := s.stats.CountByCategory(s.ctx, 100)
			s.NoError(err)

			actual := make(map[uint]int64)
			for _, c := range count {
				actual[c.CategoryId] = c.Count
			}
			expected := make(map[uint]int64)
			if test.expectedCat1Count > 0 {
				expected[1] = test.expectedCat1Count
			}
			if test.expectedCat2Count > 0 {
				expected[2] = test.expectedCat2Count
			}

			s.Equal(expected, actual)

		})
	}
}

func (s *RepoTestSuite) TestCountByCategoryAndMonth() {
	var november time.Time = time.Date(2020, 11, 4, 1, 1, 1, 0, time.UTC)
	var february time.Time = time.Date(2021, 2, 4, 1, 1, 1, 0, time.UTC)
	var cat1 = entity.Category{ID: 1, Type: entity.TypeCategory, Name: "Cat1"}
	var cat2 = entity.Category{ID: 2, Type: entity.TypeCategory, Name: "Cat2"}
	tests := map[string]struct {
		dreams           []entity.Dream
		exepectedNovCat1 int64
		exepectedNovCat2 int64
		exepectedFebCat1 int64
		exepectedFebCat2 int64
	}{
		"one dream, one cat": {
			exepectedNovCat1: 1,
			dreams:           entity.Dreams{{Date: november, Categories: entity.Categories{cat1}}},
		},
		"one dream, two cats": {
			exepectedNovCat1: 1, exepectedNovCat2: 1,
			dreams: entity.Dreams{{Date: november, Categories: entity.Categories{cat1, cat2}}},
		},
		"two dreams, distinct cats": {
			exepectedNovCat1: 1, exepectedNovCat2: 1,
			dreams: entity.Dreams{
				entity.Dream{Date: november, Categories: entity.Categories{cat1}},
				entity.Dream{Date: november, Categories: entity.Categories{cat2}},
			}},
		"two dreams, shared cats": {
			exepectedNovCat1: 1, exepectedNovCat2: 2,
			dreams: entity.Dreams{
				entity.Dream{Date: november, Categories: entity.Categories{cat1, cat2}},
				entity.Dream{Date: november, Categories: entity.Categories{cat2}},
			}},
		"different months": {
			exepectedNovCat2: 1, exepectedFebCat1: 1, exepectedFebCat2: 1,
			dreams: entity.Dreams{
				entity.Dream{Date: february, Categories: entity.Categories{cat1, cat2}},
				entity.Dream{Date: november, Categories: entity.Categories{cat2}},
			}},
	}
	for name, test := range tests {
		s.Run(name, func() {
			gorm.G[entity.Dream](s.db).CreateInBatches(s.ctx, &test.dreams, 100)
			aggs, err := s.stats.CountByCategoryAndMonth(s.ctx)
			s.NoError(err)

			actual := make(map[string]map[uint]int64)
			for _, agg := range aggs {
				if actual[agg.Month] == nil {
					actual[agg.Month] = make(map[uint]int64)
				}
				actual[agg.Month][agg.CategoryId] = agg.Count
			}
			expected := make(map[string]map[uint]int64)
			if test.exepectedFebCat1 > 0 {
				if expected["2021/02"] == nil {
					expected["2021/02"] = make(map[uint]int64)
				}
				expected["2021/02"][1] = test.exepectedFebCat1
			}
			if test.exepectedFebCat2 > 0 {
				if expected["2021/02"] == nil {
					expected["2021/02"] = make(map[uint]int64)
				}
				expected["2021/02"][2] = test.exepectedFebCat2
			}
			if test.exepectedNovCat1 > 0 {
				if expected["2020/11"] == nil {
					expected["2020/11"] = make(map[uint]int64)
				}
				expected["2020/11"][1] = test.exepectedNovCat1
			}
			if test.exepectedNovCat2 > 0 {
				if expected["2020/11"] == nil {
					expected["2020/11"] = make(map[uint]int64)
				}
				expected["2020/11"][2] = test.exepectedNovCat2
			}
			s.Equal(expected, actual)
		})
	}

}

func (s *RepoTestSuite) TestCountByDreamAndMonth() {
	tests := map[string]struct {
		numberOfDreamsInNov int64
		numberOfDreamsInFeb int64
		expectedLength      int
	}{
		"0 dreams":                  {},
		"1 dream in Nov":            {numberOfDreamsInNov: 1, expectedLength: 1},
		"5 dreams in Nov":           {numberOfDreamsInNov: 5, expectedLength: 1},
		"3 dreams in Nov, 2 in Feb": {numberOfDreamsInNov: 3, numberOfDreamsInFeb: 2, expectedLength: 2},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var dreams = []entity.Dream{}
			for range test.numberOfDreamsInNov {
				dreams = append(dreams, entity.Dream{Date: time.Date(2022, 11, 4, 0, 0, 0, 0, time.UTC),
					Categories: entity.Categories{entity.Category{ID: 1, Type: entity.TypeCategory}}})
			}
			for range test.numberOfDreamsInFeb {
				dreams = append(dreams, entity.Dream{Date: time.Date(2022, 2, 4, 0, 0, 0, 0, time.UTC)})
			}
			err := gorm.G[entity.Dream](s.db).CreateInBatches(s.ctx, &dreams, 100)
			s.Require().NoError(err)

			aggs, err := s.stats.CountByDreamAndMonth(s.ctx)
			s.NoError(err)
			if test.expectedLength == 0 {
				s.Len(aggs, 0)
				return
			}
			s.Len(aggs, test.expectedLength)

			actual := make(map[string]int64)
			for _, agg := range aggs {
				actual[agg.Month] = agg.Count
			}
			expected := map[string]int64{}
			if test.numberOfDreamsInFeb > 0 {
				expected["2022/02"] = test.numberOfDreamsInFeb
			}
			if test.numberOfDreamsInNov > 0 {
				expected["2022/11"] = test.numberOfDreamsInNov
			}
			s.Equal(expected, actual)
		})
	}

}
