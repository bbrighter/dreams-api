package usecase

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entity"
	"github.com/gin-gonic/gin"
)

type (
	Auth interface {
		SetAuth(name, pw string) string
		IsValidToken(token string) bool
		ClearToken()
		AuthMiddleware() gin.HandlerFunc
	}
)

type (
	IDreamsRepo interface {
		List(context.Context, bool) (entity.Dreams, error)
		Get(context.Context, uint, bool) (entity.Dream, error)
		Create(context.Context, *entity.Dream) (uint, error)
		Update(context.Context, uint, map[string]any) error
		Delete(context.Context, uint) error
		ToggleVisibility(context.Context, uint) error
	}

	ICategoriesRepo interface {
		List(context.Context) (entity.Categories, error)
		FindByName(ctx context.Context, catName string, catType entity.CategoryType) (uint, error)
		Create(ctx context.Context, catName string, catType entity.CategoryType) (uint, error)
		AddToDream(ctx context.Context, dreamId uint, catId uint) error
		RemoveFromDream(ctx context.Context, dreamID uint, catId uint) error
	}

	IManagementRepo interface {
		Update(context.Context, uint, map[string]any) error
		Delete(context.Context, uint) error
		Merge(context.Context, uint, uint, string) error
	}

	IStatisticsRepo interface {
		CountByCategory(ctx context.Context, limit int) (entity.CategoriesCount, error)
		CountByDreamAndMonth(ctx context.Context) ([]entity.CountByMonth, error)
		CountByCategoryAndMonth(ctx context.Context) ([]entity.CountByCatAndMonth, error)
	}

	IAuthRepo interface {
		SetAuth(name, pw string) string
		IsValidToken(token string) bool
		ClearToken()
	}
)
