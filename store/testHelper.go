package store

import (
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
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
		var typeTables []interface{} = []interface{}{
			Dream{},
			Person{},
			Tag{},
		}
		for _, table := range typeTables {
			err = repo.db.Migrator().DropTable(&table)
			assert.NoError(t, err)
		}
		var fixedTables []string = []string{
			"migrations",
			"tags_dreams",
			"people_dreams",
		}
		for _, table := range fixedTables {
			err = repo.db.Migrator().DropTable(table)
			assert.NoError(t, err)
		}
	}
	return repo, deferedFunc
}

func CreateTestDream(numberOfTags int, numberOfPersons int, t *testing.T) Dream {
	var tags []Tag
	i := 0
	for i < numberOfTags {
		tag := Tag{Title: randomdata.Noun()}
		tags = append(tags, tag)
		i++
	}
	var persons []Person
	j := 0
	for j < numberOfPersons {
		person := Person{Name: randomdata.FirstName(0)}
		persons = append(persons, person)
		j++
	}
	var dream Dream = Dream{
		Date:        time.Now(),
		Description: randomdata.RandStringRunes(100),
		Tags:        tags,
		Persons:     persons,
	}

	repo := initTestRepo()
	err := repo.db.Create(&dream).Error
	assert.NoError(t, err)
	return dream
}

func CreateTestTag(t *testing.T) Tag {
	var tag Tag = Tag{Title: randomdata.Noun()}
	repo := initTestRepo()
	err := repo.db.Create(&tag).Error
	assert.NoError(t, err)
	return tag
}

func AddTestPersonToDream(dream *Dream, t *testing.T) {
	repo := initTestRepo()

	var person Person = Person{Name: randomdata.FirstName(0)}
	dream.Persons = append(dream.Persons, person)
	repo.db.Save(&dream)
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
