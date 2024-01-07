package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var tags []Tag
	// Test no tags exists
	tags = repo.GetTags()
	assert.Len(t, tags, 0)

	// Test tags exists
	CreateTestTag(t)
	tags = repo.GetTags()
	assert.Len(t, tags, 1)
}

func TestAddTagToDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var err error
	CreateTestDream(t)

	_, err = repo.AddTagToDream("tag", 1)
	assert.NoError(t, err)

	_, err = repo.AddTagToDream("tag", 1)
	assert.NoError(t, err)

	_, err = repo.AddTagToDream("tag", 2)
	assert.Error(t, err)

	CreateTestDream(t)
	_, err = repo.AddTagToDream("tag", 2)
	assert.NoError(t, err)
}

func TestRemoveTagFromDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var err error
	var rowsAffected int64
	var dream Dream
	var tag Tag

	// Remove tag from dream
	dream = CreateTestTagAndDream(t)
	tag = dream.Tags[0]

	err = repo.RemoveTagFromDream(tag.ID, dream.ID)

	assert.NoError(t, err)
	rowsAffected = repo.db.First(&dream).RowsAffected
	assert.EqualValues(t, 1, rowsAffected)
	rowsAffected = repo.db.First(&tag).RowsAffected
	assert.EqualValues(t, 0, rowsAffected)
	rowsAffected = repo.db.First("tags_dreams").RowsAffected
	assert.EqualValues(t, 0, rowsAffected)

	CleanTestEntries(t)

	// Remove tag from dream, but tag still exists for other dream
	CreateTestTagAndDream(t)
	dream = CreateTestTagAndDream(t)
	tag = dream.Tags[0]

	err = repo.RemoveTagFromDream(tag.ID, dream.ID)

	assert.NoError(t, err)
	rowsAffected = repo.db.First(&dream).RowsAffected
	assert.EqualValues(t, 1, rowsAffected)
	rowsAffected = repo.db.First(&tag).RowsAffected
	assert.EqualValues(t, 1, rowsAffected)
	rowsAffected = repo.db.First("tags_dreams").RowsAffected
	assert.EqualValues(t, 0, rowsAffected)

	// Remove non-existing tag from dream
	dream = CreateTestTagAndDream(t)

	err = repo.RemoveTagFromDream(100, dream.ID)

	assert.Error(t, err)

	CleanTestEntries(t)

	// Remove tag from non-existing dream
	dream = CreateTestTagAndDream(t)
	tag = dream.Tags[0]

	err = repo.RemoveTagFromDream(tag.ID, 200)

	assert.Error(t, err)
}
