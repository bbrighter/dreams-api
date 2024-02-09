package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountCategories(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var cc []CategoryCount
	cc = repo.CountCategories()
	assert.Len(t, cc, 0)

	CreateTestDream(1, 0, t)
	cc = repo.CountCategories()
	assert.Len(t, cc, 1)
	assert.Equal(t, 1, cc[0].Count)
	assert.EqualValues(t, 1, cc[0].CategoryID)
}
