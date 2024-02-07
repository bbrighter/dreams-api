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
	var tags []Tag
	CreateTestDream(0, 0, t)

	tags, err = repo.AddTagToDream("tag", 1)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)

	tags, err = repo.AddTagToDream("tag", 1)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)

	_, err = repo.AddTagToDream("tag", 2)
	assert.Error(t, err)

	CreateTestDream(0, 0, t)
	tags, err = repo.AddTagToDream("tag", 2)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
}

func TestRemoveTagFromDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var err error
	var rowsAffected int64
	var dream Dream
	var tag Tag

	// Remove tag from dream
	dream = CreateTestDream(1, 0, t)
	tag = dream.Tags[0]

	_, err = repo.RemoveTagFromDream(tag.ID, dream.ID)

	assert.NoError(t, err)
	rowsAffected = repo.db.First(&dream).RowsAffected
	assert.EqualValues(t, 1, rowsAffected)
	rowsAffected = repo.db.First(&tag).RowsAffected
	assert.EqualValues(t, 0, rowsAffected)
	rowsAffected = repo.db.First("tags_dreams").RowsAffected
	assert.EqualValues(t, 0, rowsAffected)

	CleanTestEntries(t)

	// Remove tag from dream, but tag still exists for other dream
	// CreateTestTagAndDream(t)
	// dream = CreateTestTagAndDream(t)
	// tag = dream.Tags[0]

	// _, err = repo.RemoveTagFromDream(tag.ID, dream.ID)

	// assert.NoError(t, err)
	// rowsAffected = repo.db.First(&dream).RowsAffected
	// assert.EqualValues(t, 1, rowsAffected)
	// rowsAffected = repo.db.First(&tag).RowsAffected
	// assert.EqualValues(t, 1, rowsAffected)
	// rowsAffected = repo.db.First("tags_dreams").RowsAffected
	// assert.EqualValues(t, 0, rowsAffected)

	// Remove non-existing tag from dream
	dream = CreateTestDream(1, 0, t)

	_, err = repo.RemoveTagFromDream(100, dream.ID)

	assert.Error(t, err)

	CleanTestEntries(t)

	// Remove tag from non-existing dream
	dream = CreateTestDream(1, 0, t)
	tag = dream.Tags[0]

	_, err = repo.RemoveTagFromDream(tag.ID, 200)

	assert.Error(t, err)
}
