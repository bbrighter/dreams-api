package usecase

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"go.uber.org/zap"
)

type DreamsUseCase struct {
	BaseUseCase
	repo IDreamsRepo
}

func NewDreamUseCase(r IDreamsRepo, logger *zap.Logger) *DreamsUseCase {
	return &DreamsUseCase{
		repo: r, BaseUseCase: NewBaseUseCase(logger),
	}
}

func (uc *DreamsUseCase) List(includes []entity.Includes) entity.Dreams {
	return uc.repo.List(false, includes)
}

func (uc *DreamsUseCase) Get(id uint) (entity.Dream, error) {
	dream, err := uc.repo.Get(id, false)
	uc.HandleError(err)
	return dream, err
}

func (uc *DreamsUseCase) Create(date time.Time) (uint, error) {
	dream := entity.Dream{Date: date}
	id, err := uc.repo.Create(dream)
	uc.HandleError(err)
	return id, err
}

func (uc *DreamsUseCase) Update(id uint, date *time.Time, description *string, rating *int) error {
	updates := make(map[string]any)
	if date != nil {
		updates["date"] = *date
	}
	if description != nil {
		updates["description"] = *description
	}
	if rating != nil {
		updates["rating"] = *rating
	}
	if len(updates) == 0 {
		return entity.ErrorBadParamWithReasons("no params provided")
	}
	err := uc.repo.Update(id, updates)
	uc.HandleError(err)
	return err
}

func (uc *DreamsUseCase) Delete(id uint) (entity.Categories, error) {
	dream := entity.Dream{ID: id}
	cats, err := uc.repo.Delete(dream)
	uc.HandleError(err)
	return cats, err
}

func (uc *DreamsUseCase) Finalize(id uint) error {
	dream, err := uc.repo.Get(id, true)
	if uc.HandleError(err) {
		return err
	}
	if dream.Rating == nil {
		return entity.ErrorBadParam
	}

	updates := map[string]any{
		"finalized": true,
	}
	err = uc.repo.Update(id, updates)
	uc.HandleError(err)
	return err
}
