package controller

import (
	"net/http"
	"time"
)

func (s *apiTestSuite) TestStatistics() {
	tests := []apiTest{
		{name: "get monthly count, no data", method: http.MethodGet, url: "/count-categories/monthly",
			statusCode: http.StatusOK,
			response:   Statistics{Statistics: []Statistic{}},
		},
		{name: "add dream", method: http.MethodPost, url: "/dreams", body: PostDreamRequest{Date: time.Date(2026, 9, 5, 13, 10, 0, 0, time.UTC)},
			statusCode: http.StatusCreated,
		},
		{name: "add second dream", method: http.MethodPost, url: "/dreams", body: PostDreamRequest{Date: time.Date(2026, 10, 5, 13, 10, 0, 0, time.UTC)},
			statusCode: http.StatusCreated,
		},
		{name: "add category", method: http.MethodPost, url: "/dreams/1/categories", body: PostCategoryRequestBody{Name: "Cat", Type: "category"},
			statusCode: http.StatusCreated,
		},
		{name: "get monthly count with data", method: http.MethodGet, url: "/count-categories/monthly",
			statusCode: http.StatusOK,
			response: Statistics{Statistics: []Statistic{
				{
					Month:      "2026/09",
					Categories: []CountByCat{{CategoryId: 1, Count: 1}},
					DreamCount: 1,
				},
				{
					Month:      "2026/10",
					Categories: []CountByCat{},
					DreamCount: 1,
				}}},
		},
	}

	for _, test := range tests {
		s.evaluate(test)
	}
}
