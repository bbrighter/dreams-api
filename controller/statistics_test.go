package controller

import (
	"net/http"
	"time"
)

// func TestGetStatistics(t *testing.T) {
// 	g := setupApiTest(t)

// 	tests := []apiTest{
// 		{name: "get empty statistics", method: http.MethodGet, url: "/statistics?limit=40", statusCode: http.StatusOK,
// 			response: entity.CategoriesCountResponse{Categories: []entity.CountResponse{}, Persons: []entity.CountResponse{}}},
// 		{name: "setup dream", method: http.MethodPost, url: "/dreams", statusCode: http.StatusCreated,
// 			body: PostDreamRequest{Date: time.Now()}},
// 		{name: "setup person", method: http.MethodPut, url: "/dreams/1/persons?name=Pers", statusCode: http.StatusOK},
// 		{name: "setup category", method: http.MethodPut, url: "/dreams/1/categories?name=Cat", statusCode: http.StatusOK},
// 		{name: "get filled statistics", method: http.MethodGet, url: "/statistics?limit=40", statusCode: http.StatusOK,
// 			response: entity.CategoriesCountResponse{
// 				Categories: []entity.CountResponse{{ID: 2, Count: 1}},
// 				Persons:    []entity.CountResponse{{ID: 1, Count: 1}},
// 			}},
// 	}

// 	for _, test := range tests {
// 		test.evaluate(t, g)
// 	}
// }

func (s *apiTestSuite) TestStatistics() {
	tests := []apiTest{
		{name: "get monthly count, no data", method: http.MethodGet, url: "/count-categories/monthly",
			statusCode: http.StatusOK,
			response:   Statistics{Statistics: []Statistic{}},
		},
		{name: "add dream", method: http.MethodPost, url: "/dreams", body: PostDreamRequest{Date: time.Date(2026, 9, 5, 13, 10, 0, 0, time.UTC)},
			statusCode: http.StatusCreated,
		},
		{name: "add category", method: http.MethodPost, url: "/dreams/1/categories", body: PostCategoryRequestBody{Name: "Cat", Type: "category"},
			statusCode: http.StatusCreated,
		},
		{name: "get monthly count with data", method: http.MethodGet, url: "/count-categories/monthly",
			statusCode: http.StatusOK,
			response: Statistics{Statistics: []Statistic{{
				Month:      "2026/09",
				Categories: []CountByCat{{CategoryId: 1, Count: 1}},
				DreamCount: 1,
			}}},
		},
	}

	for _, test := range tests {
		s.evaluate(test)
	}
}
