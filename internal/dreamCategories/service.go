package dreamcategories

import (
	"context"
	"strings"

	"github.com/bbrighter/dreams-api/internal/categories"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"github.com/bbrighter/dreams-api/internal/entities"
	unitofwork "github.com/bbrighter/dreams-api/internal/unitOfWork"
)

type DreamCategoriesService struct {
	d   *dreams.DreamsRepo
	c   *categories.CategoriesRepo
	uow *unitofwork.UnitOfWork
}

func NewDreamCategoriesService(
	d *dreams.DreamsRepo,
	c *categories.CategoriesRepo,
	uow *unitofwork.UnitOfWork,
) *DreamCategoriesService {
	return &DreamCategoriesService{d: d, c: c, uow: uow}
}

func (s *DreamCategoriesService) AddNewCategoryToDream(
	ctx context.Context,
	dreamId uint,
	catName string,
	catType entities.CategoryType,
) (uint, error) {
	name := strings.TrimSpace(catName)
	var cat = &entities.Category{Name: name, Type: catType}

	err := s.uow.Do(ctx,
		func(dr *dreams.DreamsRepo, cr *categories.CategoriesRepo) error {
			if err := cr.CreateCategory(ctx, cat); err != nil {
				return err
			}
			if err := dr.AddCategoryToDream(ctx, dreamId, cat); err != nil {
				return err
			}
			return nil
		})
	return cat.ID, err
}

func (s *DreamCategoriesService) AddCategoryToDream(ctx context.Context, dreamId uint, catId uint) error {
	return s.d.AddCategoryToDream(ctx, dreamId, &entities.Category{ID: catId})
}

func (s *DreamCategoriesService) RemoveCategoryToDream(ctx context.Context, dreamId uint, catId uint) error {
	return s.d.RemoveCategoryFromDream(ctx, dreamId, catId)
}
