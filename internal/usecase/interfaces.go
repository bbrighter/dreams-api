package usecase

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/gin-gonic/gin"
)

type (
	Dreams interface {
		List([]entity.Includes) entity.Dreams
		Get(uint) (entity.Dream, error)
		Create(time.Time) (uint, error)
		Update(uint, *time.Time, *string, *int) error
		Delete(uint) (entity.Categories, error)
		Finalize(uint) error
	}

	PrivateDreams interface {
		ToggleVisibility(uint) error
		List() entity.Dreams
		Get(uint) (entity.Dream, error)
	}

	CategoriesLister interface {
		List([]entity.Includes) entity.Categories
	}

	CategoriesAdderRemover interface {
		AddCategoryToDream(categoryName string, dreamId uint) (entity.Categories, error)
		AddPersonToDream(categoryName string, dreamId uint) (entity.Categories, error)
		RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error)
	}

	ICategoriesManager interface {
		ChangeName(uint, string) error
		ChangeType(uint, entity.CategoryType) error
		Delete(uint) error
		Merge(uint, uint, string) (entity.Categories, error)
	}

	Statistics interface {
		GetStatistics(includeHidden bool, limit int) (entity.Counts, entity.Counts)
	}

	Auth interface {
		SetAuth(name, pw string) string
		IsValidToken(token string) bool
		ClearToken()
		AuthMiddleware() gin.HandlerFunc
	}
)

type (
	IDreamsRepo interface {
		List(bool, []entity.Includes) entity.Dreams
		Get(uint, bool) (entity.Dream, error)
		Create(entity.Dream) (uint, error)
		Update(uint, map[string]any) error
		Delete(entity.Dream) (entity.Categories, error)
	}

	ICategoriesRepo interface {
		List([]entity.Includes) entity.Categories
		AddToDream(string, entity.Dream, entity.CategoryType) error
		RemoveFromDream(entity.Category, entity.Dream) error
		Update(uint, map[string]any) error
		Delete(uint) error
		CountByNameAndType(string, entity.CategoryType) int64
		First(uint) (entity.Category, error)
		Merge(uint, uint, string) error
	}

	IStatisticsRepo interface {
		CountCategories(bool, int) (entity.Counts, entity.Counts)
	}

	IAuthRepo interface {
		SetAuth(name, pw string) string
		IsValidToken(token string) bool
		ClearToken()
	}
)
