package entity_test

import (
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestDreamToResponse(t *testing.T) {
	t.Parallel()

	var dream = entity.Dream{
		ID:          1,
		Date:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "desc",
		Visible:     true,
		Finalized:   true,
		Categories: []entity.Category{
			{ID: 1, Name: "Category", Type: entity.TypeCategory, Dreams: []entity.Dream{{ID: 1}}},
			{ID: 2, Name: "Person", Type: entity.TypePerson, Dreams: []entity.Dream{{ID: 1}}},
		},
	}

	var resp entity.DreamResponse = dream.ToResponse()
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, "desc", resp.Description)
	assert.Equal(t, true, resp.Visible)
	assert.True(t, resp.Finalized)
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), resp.Date)
	assert.Len(t, resp.Categories, 1)
	assert.Equal(t, "Category", resp.Categories[0].Name)
	assert.EqualValues(t, 1, resp.Categories[0].ID)
	assert.Len(t, resp.Persons, 1)
	assert.EqualValues(t, 2, resp.Persons[0].ID)
	assert.Equal(t, "Person", resp.Persons[0].Name)

}

func TestDreamsToDreamsResponse(t *testing.T) {
	t.Parallel()

	var dreams entity.Dreams

	// Empty input gives {dreams: []}
	var emptyResp entity.DreamsResponse = dreams.ToResponse()
	assert.NotNil(t, emptyResp.Dreams)
	assert.Len(t, emptyResp.Dreams, 0)

	// Non-empty input gives {dreams: [...]}
	dreams = entity.Dreams{
		{
			ID:          1,
			Date:        time.Now(),
			Description: "desc",
			Categories:  []entity.Category{{ID: 1, Name: "Category"}, {ID: 2, Name: "Person"}},
			Visible:     true,
			Finalized:   true,
		},
	}
	var resp entity.DreamsResponse = dreams.ToResponse()
	assert.Len(t, resp.Dreams, 1)
}
