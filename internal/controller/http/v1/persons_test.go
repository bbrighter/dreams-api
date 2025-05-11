package v1

import (
	"net/http/httptest"
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/bbrighter/dreams-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testUseCasePersons struct{}

func (tu testUseCasePersons) List() entity.Persons {
	return entity.Persons{entity.Person{
		ID: 1, Name: "Name", Dreams: []entity.Dream{{ID: 10}},
	}}
}

func TestGetAllPersons(t *testing.T) {
	tests := []struct {
		name           string
		uc             usecase.PersonsLister
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "ok",
			uc:             testUseCasePersons{},
			expectedStatus: 200,
			expectedBody:   `{"persons":[{"id":1,"name":"Name"}]}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			r := &personsRoute{p: test.uc}

			r.GetAll(c)

			assert.Equal(t, test.expectedStatus, rec.Code)
			assert.Equal(t, test.expectedBody, rec.Body.String())
		})
	}
}
