package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllPersons(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	CreateTestDream(0, 10, t)

	persons := repo.GetAllPersons()

	assert.Len(t, persons, 10)
}

func TestAddPersonToDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	var err error
	dream := CreateTestDream(0, 0, t)

	_, err = repo.AddPersonToDream("new name", dream.ID)
	assert.NoError(t, err)

	// Adding same person twice does not add a person
	_, err = repo.AddPersonToDream("new name", dream.ID)
	assert.NoError(t, err)
	var persons []Person
	repo.db.Find(&persons)
	assert.Len(t, persons, 1)

	// Adding same person to another dream does not add a person
	secondDream := CreateTestDream(0, 0, t)
	_, err = repo.AddPersonToDream("new name", secondDream.ID)
	assert.NoError(t, err)
	repo.db.Find(&persons)
	assert.Len(t, persons, 1)
}

func TestAddPersonToDreamErrors(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)

	// No dream exists
	_, err := repo.AddPersonToDream("new name", 1)
	assert.ErrorIs(t, err, ErrorNotFound)
}

func TestRemovePersonFromDream(t *testing.T) {
	repo, teardown := SetupTest(t)
	defer teardown(t)
	var dream Dream
	var err error

	// Remove person from dream
	dream = CreateTestDream(0, 1, t)
	_, err = repo.RemovePersonFromDream(dream.Persons[0].ID, dream.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, repo.db.Find(&Person{}).RowsAffected, 0)

	// Remove person from dream without person
	dream = CreateTestDream(0, 0, t)
	_, err = repo.RemovePersonFromDream(1, dream.ID)
	assert.ErrorIs(t, err, ErrorNotFound)
}
