package repository

import (
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func setupCategoriesTest(t *testing.T) *CategoriesRepo {
	repo := NewCategoriesRepo(":memory:")
	err := repo.Repo.AutoMigrate(
		&entity.Dream{},
		&entity.Category{},
		&entity.Person{},
	)
	assert.NoError(t, err)

	return repo
}

func TestGetAllCategories(t *testing.T) {
	r := setupCategoriesTest(t)

	cats := r.GetAll()
	assert.Len(t, cats, 0)

	r.Repo.Create(&entity.Category{ID: 1})
	cats = r.GetAll()
	assert.Len(t, cats, 1)
}

func TestAddCategoryToDream(t *testing.T) {
	r := setupCategoriesTest(t)
	var err error

	var dream = entity.Dream{ID: 1}
	_, err = r.AddToDream("name", dream)
	assert.Error(t, err)

	r.Repo.Create(&dream)

	var cats entity.Categories
	cats, err = r.AddToDream("name", dream)

	assert.NoError(t, err)
	assert.Len(t, cats, 1)

	cats, err = r.AddToDream("name", dream)

	assert.NoError(t, err)
	assert.Len(t, cats, 1)

	cats, err = r.AddToDream("new name", dream)

	assert.NoError(t, err)
	assert.Len(t, cats, 2)
}

func TestRemoveCategoryFromDream(t *testing.T) {
	r := setupCategoriesTest(t)
	var err error

	// No dream and no category
	var dream = entity.Dream{ID: 1}
	var cat = entity.Category{ID: 10}
	_, err = r.RemoveFromDream(cat, dream)
	assert.Error(t, err)

	// No category
	r.Repo.Create(&dream)

	_, err = r.RemoveFromDream(cat, dream)
	assert.Error(t, err)

	// Category and dream exist
	r.Repo.Create(&entity.Category{ID: 10, Dreams: []entity.Dream{dream}})

	var cats entity.Categories
	cats, err = r.RemoveFromDream(cat, dream)
	assert.NoError(t, err)
	assert.Len(t, cats, 0)

	// Category exists and cannot be removed
	var dream2 = entity.Dream{ID: 2}
	r.Repo.Create(&entity.Dreams{dream, dream2})
	r.Repo.Create(&entity.Category{ID: 10, Dreams: []entity.Dream{dream, dream2}})
	cats, err = r.RemoveFromDream(cat, dream)
	assert.NoError(t, err)
	assert.Len(t, cats, 1)
}
