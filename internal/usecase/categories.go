package usecase

import "github.com/bbrighter/dreams-api/internal/entity"

type CategoriesUseCase struct {
	repo ICategoriesRepo
}

func NewCategoriesUseCase(r ICategoriesRepo) *CategoriesUseCase {
	return &CategoriesUseCase{repo: r}
}

func (u *CategoriesUseCase) List() entity.Categories {
	return u.repo.List()
}

func (u *CategoriesUseCase) AddToDream(categoryName string, dreamId uint) (entity.Categories, error) {
	var dream = entity.Dream{ID: dreamId}
	return u.repo.AddToDream(categoryName, dream)
}

func (u *CategoriesUseCase) RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error) {
	var dream = entity.Dream{ID: dreamId}
	var category = entity.Category{ID: categoryId}
	return u.repo.RemoveFromDream(category, dream)
}
