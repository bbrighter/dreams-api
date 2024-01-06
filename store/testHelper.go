package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initTestRepo() Repo {
	return InitRepo("test.sqlite")
}

func SetupTest(t *testing.T) (Repo, func(t *testing.T)) {
	repo := initTestRepo()
	err := Migration(repo)
	assert.NoError(t, err)

	deferedFunc := func(t *testing.T) {
		var err error
		err = repo.db.Migrator().DropTable("dreams")
		assert.NoError(t, err)
		err = repo.db.Migrator().DropTable("tags")
		assert.NoError(t, err)
		err = repo.db.Migrator().DropTable("tags_dreams")
		assert.NoError(t, err)
		err = repo.db.Migrator().DropTable("migrations")
		assert.NoError(t, err)
	}
	return repo, deferedFunc
}

func CreateTestDream(t *testing.T) Dream {
	var dream Dream = Dream{Date: time.Now(), Description: "desc"}
	repo := initTestRepo()
	err := repo.db.Create(&dream).Error
	assert.NoError(t, err, "Creation of example dream failed")
	return dream
}

func CreateTestTag(t *testing.T) Tag {
	var tag Tag = Tag{Title: "Tag"}
	repo := initTestRepo()
	err := repo.db.Create(&tag).Error
	assert.NoError(t, err)
	return tag
}

func CreateTestTagAndDream(t *testing.T) Dream {
	repo := initTestRepo()
	dream := Dream{
		Date:        time.Now(),
		Description: "desc with tag",
		Tags:        []Tag{{Title: "Tag"}},
	}
	err := repo.db.Create(&dream).Error
	assert.NoError(t, err)
	return dream
}

func CleanTestEntries(t *testing.T) {
	repo := initTestRepo()
	tx := repo.db.Session(&gorm.Session{AllowGlobalUpdate: true})
	err := tx.Delete(&Dream{}).Error
	assert.NoError(t, err)
	err = tx.Delete(&Tag{}).Error
	assert.NoError(t, err)
	err = tx.Table("tags_dreams").Delete("tags_dreams").Error
	assert.NoError(t, err)

}
