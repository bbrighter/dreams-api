package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountCategories(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var cc []Count
	cc = repo.CountCategories(false)
	assert.Len(t, cc, 0)

	visibleDream := CreateTestDream(1, 0, true, t)
	CreateTestDream(1, 0, false, t)
	cc = repo.CountCategories(false)
	assert.Len(t, cc, 1)
	assert.Equal(t, 1, cc[0].Count)
	assert.EqualValues(t, visibleDream.Categories[0].ID, cc[0].ID)

	cc = repo.CountCategories(true)
	assert.Len(t, cc, 2)
}

func TestCountPersons(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var cc []Count
	cc = repo.CountPersons(false)
	assert.Len(t, cc, 0)

	CreateTestDream(0, 12, true, t)
	cc = repo.CountPersons(false)
	assert.Len(t, cc, 12)
	assert.GreaterOrEqual(t, 1, cc[0].Count)
	assert.GreaterOrEqual(t, 1, int(cc[0].ID))
}
