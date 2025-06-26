package usecase

import (
	"strings"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type CategoriesUseCase struct {
	repo ICategoriesRepo
}

func NewCategoriesUseCase(r ICategoriesRepo) *CategoriesUseCase {
	return &CategoriesUseCase{repo: r}
}

func (u *CategoriesUseCase) List(includes []entity.Includes) entity.Categories {
	return u.repo.List(includes)
}

func (u *CategoriesUseCase) AddCategoryToDream(categoryName string, dreamId uint) (entity.Categories, error) {
	var dream = entity.Dream{ID: dreamId}
	return u.repo.AddToDream(categoryName, dream, entity.TypeCategory)
}

func (u *CategoriesUseCase) AddPersonToDream(categoryName string, dreamId uint) (entity.Categories, error) {
	var dream = entity.Dream{ID: dreamId}
	return u.repo.AddToDream(categoryName, dream, entity.TypePerson)
}

func (u *CategoriesUseCase) RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error) {
	var dream = entity.Dream{ID: dreamId}
	var category = entity.Category{ID: categoryId}
	return u.repo.RemoveFromDream(category, dream)
}

type CategoriesManager struct {
	repo ICategoriesRepo
}

func NewCategoriesManager(repo ICategoriesRepo) *CategoriesManager {
	return &CategoriesManager{repo: repo}
}

func (u *CategoriesManager) ChangeType(categoryId uint, newType entity.CategoryType) error {
	updates := map[string]any{"type": newType}
	return u.repo.Update(categoryId, updates)
}

func (u *CategoriesManager) ChangeName(categoryId uint, newName string) error {
	trimmedName := strings.TrimSpace(newName)
	cat, err := u.repo.First(categoryId)
	if err != nil {
		return err
	}
	if count := u.repo.CountByNameAndType(trimmedName, cat.Type); count > 0 {
		return entity.ErrorBadParamWithReasons("name already exists")
	}
	updates := map[string]any{"name": trimmedName}
	return u.repo.Update(categoryId, updates)
}

func (u *CategoriesManager) Delete(categoryId uint) error {
	return u.repo.Delete(categoryId)
}

func (u *CategoriesManager) Merge(sourceCategoryId uint, targetCategoryId uint, newName string) (entity.Categories, error) {
	if err := u.repo.Merge(sourceCategoryId, targetCategoryId, newName); err != nil {
		return entity.Categories{}, err
	}
	return u.repo.List([]entity.Includes{}), nil
}
