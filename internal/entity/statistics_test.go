package entity_test

import (
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestStatisticsToResponse(t *testing.T) {
	t.Parallel()
	var counts = entity.Counts{
		entity.Count{ID: 1, Count: 10},
		entity.Count{ID: 2, Count: 2},
	}

	resp := counts.ToResponse()
	assert.Len(t, resp, 2)
	assert.Equal(t, resp, []entity.CountResponse{
		{ID: 1, Count: 10},
		{ID: 2, Count: 2},
	})
}
