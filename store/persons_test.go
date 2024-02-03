package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPersonsForDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)
	var persons []Person

	persons = repo.GetPersonsForDream(1)
	assert.Len(t, persons, 0)

	dream := CreateTestDream(0, 1, t)

	persons = repo.GetPersonsForDream(dream.ID)
	assert.Len(t, persons, 1)
}

func TestAddPersonToDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var err error
	dream := CreateTestDream(0, 0, t)

	err = repo.AddPersonToDream("new name", dream.ID)
	assert.NoError(t, err)

	// Adding same person twice does not add a person
	err = repo.AddPersonToDream("new name", dream.ID)
	assert.NoError(t, err)
	var persons []Person
	repo.db.Find(&persons)
	assert.Len(t, persons, 1)

	// Adding same person to another dream does not add a person
	secondDream := CreateTestDream(0, 0, t)
	err = repo.AddPersonToDream("new name", secondDream.ID)
	assert.NoError(t, err)
	repo.db.Find(&persons)
	assert.Len(t, persons, 1)
}

func TestAddPersonToDreamErrors(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	// No dream exists
	err := repo.AddPersonToDream("new name", 1)
	assert.ErrorIs(t, err, ErrorNotFound)
}

func TestRemovePersonFromDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)
	var dream Dream
	var err error

	// Remove person from dream
	dream = CreateTestDream(0, 1, t)
	err = repo.RemovePersonFromDream(dream.Persons[0].ID, dream.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, repo.db.Find(&Person{}).RowsAffected, 0)

	// Remove person from dream without person
	dream = CreateTestDream(0, 0, t)
	err = repo.RemovePersonFromDream(1, dream.ID)
	assert.ErrorIs(t, err, ErrorNotFound)
}
