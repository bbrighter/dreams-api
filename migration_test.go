package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (Repo, func(t *testing.T)) {
	repo := InitRepo("test.sqlite")
	err := Migration(repo.db)
	assert.NoError(t, err)

	deferedFunc := func(t *testing.T) {
		// assert.NoError(t, RollbackTo(repo.db, initialMigration))
		// assert.NoError(t, Rollback(repo.db))
		var err error
		err = repo.db.Migrator().DropTable("dreams")
		assert.NoError(t, err)
		err = repo.db.Migrator().DropTable("tags")
		assert.NoError(t, err)
		err = repo.db.Migrator().DropTable("migrations")
		assert.NoError(t, err)
	}
	return repo, deferedFunc
}

func TestMigration(t *testing.T) {
	_, setup := setupTest(t)
	defer setup(t)
}
