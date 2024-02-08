package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCategories(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var categories []Category
	// Test categories exists
	categories = repo.GetCategories()
	assert.Len(t, categories, 0)

	// Test categories exists
	CreateTestCategory(t)
	categories = repo.GetCategories()
	assert.Len(t, categories, 1)
}

func TestAddCategoryToDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var err error
	var categories []Category
	CreateTestDream(0, 0, t)

	categories, err = repo.AddCategoryToDream("category", 1)
	assert.NoError(t, err)
	assert.Len(t, categories, 1)

	categories, err = repo.AddCategoryToDream("category", 1)
	assert.NoError(t, err)
	assert.Len(t, categories, 1)

	_, err = repo.AddCategoryToDream("category", 2)
	assert.Error(t, err)

	CreateTestDream(0, 0, t)
	categories, err = repo.AddCategoryToDream("category", 2)
	assert.NoError(t, err)
	assert.Len(t, categories, 1)
}

func TestRemoveCategoryFromDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var err error
	var rowsAffected int64
	var dream Dream
	var category Category

	// Remove category from dream
	dream = CreateTestDream(1, 0, t)
	category = dream.Categories[0]

	_, err = repo.RemoveCategoryFromDream(category.ID, dream.ID)

	assert.NoError(t, err)
	rowsAffected = repo.db.First(&dream).RowsAffected
	assert.EqualValues(t, 1, rowsAffected)
	rowsAffected = repo.db.First(&category).RowsAffected
	assert.EqualValues(t, 0, rowsAffected)
	rowsAffected = repo.db.First("categories_dreams").RowsAffected
	assert.EqualValues(t, 0, rowsAffected)

	CleanTestEntries(t)

	dream = CreateTestDream(1, 0, t)

	_, err = repo.RemoveCategoryFromDream(100, dream.ID)

	assert.Error(t, err)

	CleanTestEntries(t)

	// Remove category from non-existing dream
	dream = CreateTestDream(1, 0, t)
	category = dream.Categories[0]

	_, err = repo.RemoveCategoryFromDream(category.ID, 200)

	assert.Error(t, err)
}
