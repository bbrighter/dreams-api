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

	// Test no dreams exist
	dreams = repo.GetDreams()

	assert.Len(t, dreams, 0)

	// Test one dream exists
	CreateTestDream(t)

	dreams = repo.GetDreams()

	assert.Len(t, dreams, 1)

	// Test one dream and tag exists
	CreateTestTagAndDream(t)

	dreams = repo.GetDreams()

	assert.Len(t, dreams, 2)
	dream := dreams[1]
	assert.Len(t, dream.Tags, 1)
	tag := dream.Tags[0]
	assert.Equal(t, "Tag", tag.Title)
}

func TestGetDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var dream Dream
	var err error

	// Dream doesn't exist
	_, err = repo.GetDream(1)
	assert.EqualError(t, err, "Not found")

	// Dream exists
	CreateTestDream(t)

	dream, err = repo.GetDream(1)
	assert.NoError(t, err)
	assert.EqualValues(t, dream.Description, "desc")

	// Dream exists and has tag
	CreateTestTagAndDream(t)

	dream, err = repo.GetDream(2)
	assert.NoError(t, err)
	assert.Len(t, dream.Tags, 1)

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
	CreateTestDream(t)

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
	assert.EqualError(t, err, "Not found")

	CreateTestDream(t)

	err = repo.DeleteDream(1)
	assert.NoError(t, err)
}

func TestDeleteDreamAlsoDeletesTag(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)
	dream := CreateTestTagAndDream(t)

	err := repo.DeleteDream(dream.ID)

	assert.NoError(t, err)
	var tags []Tag
	repo.db.Find(&tags)
	assert.Len(t, tags, 0)

}
