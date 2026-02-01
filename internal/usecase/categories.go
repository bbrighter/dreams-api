package usecase

import (
	"context"
	"strings"

	"github.com/bbrighter/dreams-api/internal/entity"
	"go.uber.org/zap"
)

type CategoriesUseCase struct {
	BaseUseCase
	repo ICategoriesRepo
}

type ICategories interface {
	List(context.Context) (entity.Categories, error)
	AddNewCategoryToDream(ctx context.Context, dreamId uint, catName string, catType entity.CategoryType) (uint, error)
	AddCategoryToDream(ctx context.Context, dreamId uint, catId uint) error
	RemoveCategoryFromDream(ctx context.Context, dreamID uint, catId uint) error
}

func NewCategoriesUseCase(r ICategoriesRepo, logger *zap.Logger) *CategoriesUseCase {
	return &CategoriesUseCase{repo: r, BaseUseCase: NewBaseUseCase(logger)}
}

func (u *CategoriesUseCase) List(ctx context.Context) (entity.Categories, error) {
	return u.repo.List(ctx)
}

func (u *CategoriesUseCase) AddNewCategoryToDream(ctx context.Context, dreamId uint, name string, catType entity.CategoryType) (uint, error) {
	name = strings.TrimSpace(name)
	catId, err := u.repo.Create(ctx, name, catType)
	if err != nil {
		return catId, err
	}
	err = u.repo.AddToDream(ctx, dreamId, catId)
	return catId, err
}

func (u *CategoriesUseCase) AddCategoryToDream(ctx context.Context, dreamId uint, catId uint) error {
	return u.repo.AddToDream(ctx, dreamId, catId)
}

func (u *CategoriesUseCase) RemoveCategoryFromDream(ctx context.Context, dreamId uint, catId uint) error {
	return u.repo.RemoveFromDream(ctx, dreamId, catId)
}
