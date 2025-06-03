package usecase

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type (
	Dreams interface {
		List() entity.Dreams
		Get(uint) (entity.Dream, error)
		Create(time.Time) (uint, error)
		Update(uint, time.Time, string) error
		Delete(uint) (entity.Categories, entity.Persons, error)
		Finalize(uint) error
	}

	PrivateDreams interface {
		ToggleVisibility(uint) error
		List() entity.Dreams
		Get(uint) (entity.Dream, error)
	}

	CategoriesLister interface {
		List() entity.Categories
	}

	CategoriesAdderRemover interface {
		AddToDream(categoryName string, dreamId uint) (entity.Categories, error)
		RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error)
	}

	PersonsLister interface {
		List() entity.Persons
	}

	PersonsAdderRemover interface {
		AddToDream(personName string, dreamId uint) (entity.Persons, error)
		RemoveFromDream(personId uint, dreamId uint) (entity.Persons, error)
	}

	Statistics interface {
		GetStatistics(includeHidden bool, limit int) (entity.Counts, entity.Counts)
	}
)

type (
	IDreamsRepo interface {
		List(bool) entity.Dreams
		Get(uint, bool) (entity.Dream, error)
		Create(entity.Dream) (uint, error)
		Update(entity.Dream) error
		Delete(entity.Dream) (entity.Categories, entity.Persons, error)
		ToggleVisibility(entity.Dream) error
		Finalize(uint) error
	}

	ICategoriesRepo interface {
		List() entity.Categories
		AddToDream(string, entity.Dream) (entity.Categories, error)
		RemoveFromDream(entity.Category, entity.Dream) (entity.Categories, error)
	}

	IPersonsRepo interface {
		List() entity.Persons
		AddToDream(string, entity.Dream) (entity.Persons, error)
		RemoveFromDream(entity.Person, entity.Dream) (entity.Persons, error)
	}

	IStatisticsRepo interface {
		CountPersons(bool, int) entity.Counts
		CountCategories(bool, int) entity.Counts
	}
)
