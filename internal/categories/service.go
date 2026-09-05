package categories

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entities"
)

type CategoriesService struct {
	repo *CategoriesRepo
}

func NewCategoriesService(repo *CategoriesRepo) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) ListCategories(ctx context.Context) ([]entities.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *CategoriesService) ListAndCountCategories(ctx context.Context) ([]CountByCategory, error) {
	return s.repo.ListAndCountCategories(ctx)
}

func (s *CategoriesService) UpdateCategory(ctx context.Context, id uint, name *string, catType *entities.CategoryType) error {
	var updates = make(map[string]any)
	if name != nil {
		updates["name"] = *name
	}
	if catType != nil {
		updates["type"] = *catType
	}
	return s.repo.UpdateCategory(ctx, id, updates)
}

func (s *CategoriesService) MergeCategories(
	ctx context.Context,
	sourceId uint,
	targetId uint,
	newName string,
) error {
	return s.repo.MergeCategories(ctx, sourceId, targetId, newName)
}

func (s *CategoriesService) DeleteCategory(ctx context.Context, id uint) error {
	return s.repo.DeleteCategory(ctx, id)
}
