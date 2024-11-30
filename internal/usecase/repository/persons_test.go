package repository

import (
	"testing"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupPersonsTest(t *testing.T) *PersonsRepo {
	logger, _ := zap.NewDevelopment()
	db := NewDatabase(":memory:", logger)
	repo := NewPersonsRepo(db)
	err := repo.Repo.AutoMigrate(
		&entity.Dream{},
		&entity.Category{},
		&entity.Person{},
	)
	assert.NoError(t, err)

	return repo
}
func TestGetAllPersons(t *testing.T) {
	r := setupPersonsTest(t)

	persons := r.List()
	assert.Len(t, persons, 0)

	r.Repo.Create(&entity.Person{ID: 1})
	persons = r.List()
	assert.Len(t, persons, 1)
}

func TestAddPersonToDream(t *testing.T) {
	r := setupPersonsTest(t)

	var err error
	var persons entity.Persons

	var dream = entity.Dream{ID: 1}

	// Dream doesn't exist
	_, err = r.AddToDream("name", dream)
	assert.Error(t, err)

	// Valid
	r.Repo.Create(&dream)

	persons, err = r.AddToDream("name", dream)
	assert.NoError(t, err)
	assert.Len(t, persons, 1)

	// Idempotency
	persons, err = r.AddToDream("name", dream)
	assert.NoError(t, err)
	assert.Len(t, persons, 1)

	// Different name is added
	persons, err = r.AddToDream("new name", dream)
	assert.NoError(t, err)
	assert.Len(t, persons, 2)
}

func TestRemovePersonFromDream(t *testing.T) {
	r := setupPersonsTest(t)

	var err error
	var persons entity.Persons

	var person = entity.Person{ID: 10}
	var dream = entity.Dream{ID: 1, Persons: entity.Persons{person}}

	// Dream and person not found
	_, err = r.RemoveFromDream(person, dream)
	assert.Error(t, err)

	// Dream not found
	r.Repo.Create(&person)
	_, err = r.RemoveFromDream(person, dream)
	assert.Error(t, err)

	// Valid
	r.Repo.Create(&dream)
	persons, err = r.RemoveFromDream(person, dream)
	assert.NoError(t, err)
	assert.Len(t, persons, 0)

	// One person is left
	r.Repo.Create(&entity.Person{ID: 10, Dreams: []entity.Dream{dream}})
	r.Repo.Create(&entity.Person{ID: 11, Dreams: []entity.Dream{{ID: 2}}})
	persons, err = r.RemoveFromDream(person, dream)
	assert.NoError(t, err)
	assert.Len(t, persons, 1)
}
