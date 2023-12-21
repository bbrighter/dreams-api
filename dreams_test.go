package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestDream(t *testing.T) {
	repo := InitRepo("test.sqlite")
	err := repo.db.Create(&Dream{Date: time.Now(), Description: "desc"}).Error
	assert.NoError(t, err, "Creation of example dream failed")
}

func TestGetDreams(t *testing.T) {
	repo, setup := setupTest(t)
	defer setup(t)

	dreams, err := repo.getDreams()

	assert.NoError(t, err)
	assert.Len(t, dreams, 0)

	createTestDream(t)

	dreams, err = repo.getDreams()

	assert.NoError(t, err)
	assert.Len(t, dreams, 1)
}

func TestGetDream(t *testing.T) {
	repo, setup := setupTest(t)
	defer setup(t)

	_, err := repo.getDream(1)
	assert.EqualError(t, err, "Not found")

	createTestDream(t)

	dream, err := repo.getDream(1)
	assert.NoError(t, err)
	assert.EqualValues(t, dream.Description, "desc")

}

func TestCreateDream(t *testing.T) {
	repo, setup := setupTest(t)
	defer setup(t)

	dream := Dream{
		Date:        time.Date(2023, 12, 19, 0, 0, 0, 0, time.UTC),
		Description: "desc",
	}
	id, err := repo.createDream(dream)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)

	var result Dream
	rows := repo.db.Find(&result).RowsAffected
	assert.EqualValues(t, 1, rows)
	assert.Equal(t, "desc", result.Description)
}

func TestDeleteDream(t *testing.T) {
	repo, setup := setupTest(t)
	defer setup(t)

	err := repo.deleteDream(1)
	assert.EqualError(t, err, "Not found")

	createTestDream(t)

	err = repo.deleteDream(1)
	assert.NoError(t, err)

}
