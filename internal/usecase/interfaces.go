package usecase

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type (
	Dreams interface {
		List([]entity.Includes) entity.Dreams
		Get(uint) (entity.Dream, error)
		Create(time.Time) (uint, error)
		Update(uint, time.Time, string) error
		Delete(uint) (entity.Categories, error)
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
		AddCategoryToDream(categoryName string, dreamId uint) (entity.Categories, error)
		AddPersonToDream(categoryName string, dreamId uint) (entity.Categories, error)
		RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error)
	}

	Statistics interface {
		GetStatistics(includeHidden bool, limit int) (entity.Counts, entity.Counts)
	}
)

type (
	IDreamsRepo interface {
		List(bool, []entity.Includes) entity.Dreams
		Get(uint, bool) (entity.Dream, error)
		Create(entity.Dream) (uint, error)
		Update(entity.Dream) error
		Delete(entity.Dream) (entity.Categories, error)
		ToggleVisibility(entity.Dream) error
		Finalize(uint) error
	}

	ICategoriesRepo interface {
		List() entity.Categories
		AddToDream(string, entity.Dream, entity.CategoryType) (entity.Categories, error)
		RemoveFromDream(entity.Category, entity.Dream) (entity.Categories, error)
	}

	IStatisticsRepo interface {
		CountCategories(bool, int) (entity.Counts, entity.Counts)
	}
)
