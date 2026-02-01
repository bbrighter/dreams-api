package usecase

import (
	"context"
	"strings"

	"github.com/bbrighter/dreams-api/internal/entity"
	"go.uber.org/zap"
)

type ICategoriesManager interface {
	ChangeName(context.Context, uint, string) error
	ChangeType(context.Context, uint, entity.CategoryType) error
	Delete(context.Context, uint) error
	Merge(context.Context, uint, uint, string) (entity.CategoriesCount, error)
}

type CategoriesManager struct {
	BaseUseCase
	m IManagementRepo
	s IStatisticsRepo
}

func NewCategoriesManager(m IManagementRepo, s IStatisticsRepo, logger *zap.Logger) *CategoriesManager {
	return &CategoriesManager{m: m, s: s, BaseUseCase: NewBaseUseCase(logger)}
}

func (u *CategoriesManager) ChangeType(ctx context.Context, categoryId uint, newType entity.CategoryType) error {
	updates := map[string]any{"type": newType}
	err := u.m.Update(ctx, categoryId, updates)
	return err
}

func (u *CategoriesManager) ChangeName(ctx context.Context, categoryId uint, newName string) error {
	trimmedName := strings.TrimSpace(newName)
	updates := map[string]any{"name": trimmedName}
	return u.m.Update(ctx, categoryId, updates)
}

func (u *CategoriesManager) Delete(ctx context.Context, categoryId uint) error {
	return u.m.Delete(ctx, categoryId)
}

func (u *CategoriesManager) Merge(ctx context.Context, sourceCategoryId uint, targetCategoryId uint, newName string) (entity.CategoriesCount, error) {
	if err := u.m.Merge(ctx, sourceCategoryId, targetCategoryId, newName); err != nil {
		return entity.CategoriesCount{}, err
	}
	return u.s.CountByCategory(ctx, 0)
}
