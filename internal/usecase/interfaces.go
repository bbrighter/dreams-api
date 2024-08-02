package usecase

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type (
	Dreams interface {
		GetAll(bool) entity.Dreams
		Get(uint, bool) (entity.Dream, error)
		Create(time.Time) (uint, error)
		Update(uint, time.Time, string) error
		Delete(uint) (entity.Categories, error)
		ToggleVisibility(uint) error
	}

	DreamsRepo interface {
		GetAll(bool) entity.Dreams
		GetById(uint, bool) (entity.Dream, error)
		Create(entity.Dream) (uint, error)
		Update(entity.Dream) error
		Delete(entity.Dream) (entity.Categories, error)
		ToggleVisibility(entity.Dream) error
	}

	Categories interface {
		GetAll() entity.Categories
		AddToDream(categoryName string, dreamId uint) (entity.Categories, error)
		RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error)
	}

	CategoriesRepo interface {
		GetAll() entity.Categories
		AddToDream(string, entity.Dream) (entity.Categories, error)
		RemoveFromDream(entity.Category, entity.Dream) (entity.Categories, error)
	}

	Persons interface {
		GetAll() entity.Persons
		AddToDream(personName string, dreamId uint) (entity.Persons, error)
		RemoveFromDream(personId uint, dreamId uint) (entity.Persons, error)
	}

	PersonsRepo interface {
		GetAll() entity.Persons
		AddToDream(string, entity.Dream) (entity.Persons, error)
		RemoveFromDream(entity.Person, entity.Dream) (entity.Persons, error)
	}

	Statistics interface {
		GetStatistics(bool, int) (entity.Counts, entity.Counts)
	}

	StatisticsRepo interface {
		CountPersons(bool, int) entity.Counts
		CountCategories(bool, int) entity.Counts
	}
)
