package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDreams(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var dreams []Dream
	var showAll bool = false

	// Test no dreams exist
	dreams = repo.GetDreams(&showAll)

	assert.Len(t, dreams, 0)

	// Test one dream exists
	CreateTestDream(0, 0, true, t)

	dreams = repo.GetDreams(&showAll)

	assert.Len(t, dreams, 1)

	// Test one dream and tag exists
	CreateTestDream(1, 1, true, t)
	CreateTestDream(1, 1, false, t)

	dreams = repo.GetDreams(&showAll)

	assert.Len(t, dreams, 2)
	d := dreams[1]
	assert.Len(t, d.Categories, 0)
	assert.Len(t, d.Persons, 0)

	showAll = true
	dreams = repo.GetDreams(&showAll)
	assert.Len(t, dreams, 3)
}

func TestGetDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var showAll bool = false
	var dream Dream
	var err error

	// Dream doesn't exist
	_, err = repo.GetDream(1, &showAll)
	assert.EqualError(t, err, "not found")

	// Dream exists
	o := CreateTestDream(0, 0, true, t)

	dream, err = repo.GetDream(1, &showAll)
	assert.NoError(t, err)
	assert.EqualValues(t, dream.Description, o.Description)

	// Dream exists and has tag
	CreateTestDream(1, 0, true, t)

	dream, err = repo.GetDream(2, &showAll)
	assert.NoError(t, err)
	assert.Len(t, dream.Categories, 1)

	// Invisible dream
	dreamId := CreateTestDream(0, 0, false, t).ID

	dream, err = repo.GetDream(dreamId, &showAll)
	assert.Error(t, err)

	showAll = true
	dream, err = repo.GetDream(dreamId, &showAll)
	assert.NoError(t, err)
}

func TestCreateDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	dream := Dream{
		Date:        time.Date(2023, 12, 19, 0, 0, 0, 0, time.UTC),
		Description: "desc",
	}
	id, err := repo.CreateDream(dream)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)

	var result Dream
	rows := repo.db.Find(&result).RowsAffected
	assert.EqualValues(t, 1, rows)
	assert.Equal(t, "desc", result.Description)
}

func TestUpdateDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)
	CreateTestDream(0, 0, true, t)

	var err error

	err = repo.UpdateDream(Dream{ID: 1, Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)})
	assert.NoError(t, err)

	err = repo.UpdateDream(Dream{ID: 1, Description: "new"})
	assert.NoError(t, err)

	err = repo.UpdateDream(Dream{ID: 2})
	assert.Error(t, err)
}

func TestDeleteDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	err := repo.DeleteDream(1)
	assert.EqualError(t, err, "not found")

	CreateTestDream(0, 0, true, t)

	err = repo.DeleteDream(1)
	assert.NoError(t, err)
}

func TestDeleteDreamAlsoDeletesTag(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)
	dream := CreateTestDream(1, 0, true, t)

	err := repo.DeleteDream(dream.ID)

	assert.NoError(t, err)
	var tags []Category
	repo.db.Find(&tags)
	assert.Len(t, tags, 0)

}

func TestToggleVisibility(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)
	dream := CreateTestDream(1, 0, true, t)

	visible, err := repo.TogglePrivateDream(dream.ID)

	assert.NoError(t, err)
	assert.False(t, visible)

	visible, err = repo.TogglePrivateDream(dream.ID)
	assert.NoError(t, err)
	assert.True(t, visible)
}
