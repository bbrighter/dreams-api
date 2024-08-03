package usecase

import "github.com/bbrighter/dreams-api/internal/entity"

type CategoriesUseCase struct {
	repo CategoriesRepo
}

func NewCategoriesUseCase(r CategoriesRepo) *CategoriesUseCase {
	return &CategoriesUseCase{repo: r}
}

func (u *CategoriesUseCase) GetAll() entity.Categories {
	return u.repo.GetAll()
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
