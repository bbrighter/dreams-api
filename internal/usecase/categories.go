package usecase

import (
	"strings"

	"github.com/bbrighter/dreams-api/internal/entity"
	"go.uber.org/zap"
)

type CategoriesUseCase struct {
	BaseUseCase
	repo ICategoriesRepo
}

func NewCategoriesUseCase(r ICategoriesRepo, logger *zap.Logger) *CategoriesUseCase {
	return &CategoriesUseCase{repo: r, BaseUseCase: NewBaseUseCase(logger)}
}

func (u *CategoriesUseCase) List(includes []entity.Includes) entity.Categories {
	return u.repo.List(includes)
}

func (u *CategoriesUseCase) AddCategoryToDream(categoryName string, dreamId uint) (entity.Categories, error) {
	var dream = entity.Dream{ID: dreamId}
	if err := u.repo.AddToDream(categoryName, dream, entity.TypeCategory); err != nil {
		u.HandleError(err)
		return nil, err
	}
	cats := u.repo.List([]entity.Includes{})
	return cats, nil
}

func (u *CategoriesUseCase) AddPersonToDream(categoryName string, dreamId uint) (entity.Categories, error) {
	var dream = entity.Dream{ID: dreamId}
	if err := u.repo.AddToDream(categoryName, dream, entity.TypePerson); err != nil {
		u.HandleError(err)
		return nil, err
	}
	cats := u.repo.List([]entity.Includes{})
	return cats, nil
}

func (u *CategoriesUseCase) RemoveFromDream(categoryId uint, dreamId uint) (entity.Categories, error) {
	var dream = entity.Dream{ID: dreamId}
	var category = entity.Category{ID: categoryId}
	if err := u.repo.RemoveFromDream(category, dream); err != nil {
		u.HandleError(err)
		return nil, err
	}
	cats := u.repo.List([]entity.Includes{})
	return cats, nil
}

type CategoriesManager struct {
	BaseUseCase
	repo ICategoriesRepo
}

func NewCategoriesManager(repo ICategoriesRepo, logger *zap.Logger) *CategoriesManager {
	return &CategoriesManager{repo: repo, BaseUseCase: NewBaseUseCase(logger)}
}

func (u *CategoriesManager) ChangeType(categoryId uint, newType entity.CategoryType) error {
	updates := map[string]any{"type": newType}
	err := u.repo.Update(categoryId, updates)
	u.HandleError(err)
	return err
}

func (u *CategoriesManager) ChangeName(categoryId uint, newName string) error {
	trimmedName := strings.TrimSpace(newName)
	cat, err := u.repo.First(categoryId)
	if u.HandleError(err) {
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
		u.HandleError(err)
		return entity.Categories{}, err
	}
	return u.repo.List([]entity.Includes{entity.IncludeDreamsCount}), nil
}
