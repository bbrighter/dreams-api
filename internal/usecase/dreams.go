package usecase

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type DreamsUseCase struct {
	repo IDreamsRepo
}

func NewDreamUseCase(r IDreamsRepo) *DreamsUseCase {
	return &DreamsUseCase{
		repo: r,
	}
}

func (uc *DreamsUseCase) List(includes []entity.Includes) entity.Dreams {
	return uc.repo.List(false, includes)
}

func (uc *DreamsUseCase) Get(id uint) (entity.Dream, error) {
	return uc.repo.Get(id, false)
}

func (uc *DreamsUseCase) Create(date time.Time) (uint, error) {
	dream := entity.Dream{Date: date}
	return uc.repo.Create(dream)
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
	return uc.repo.Update(id, updates)
}

func (uc *DreamsUseCase) Delete(id uint) (entity.Categories, error) {
	dream := entity.Dream{ID: id}
	return uc.repo.Delete(dream)
}

func (uc *DreamsUseCase) Finalize(id uint) error {
	dream, err := uc.repo.Get(id, true)
	if err != nil {
		return err
	}
	if dream.Rating == nil {
		return entity.ErrorBadParam
	}

	updates := map[string]any{
		"finalized": true,
	}
	return uc.repo.Update(id, updates)
}
