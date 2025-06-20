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
	return uc.repo.Update(id, updates)
}

func (uc *DreamsUseCase) Delete(id uint) (entity.Categories, error) {
	dream := entity.Dream{ID: id}
	return uc.repo.Delete(dream)
}

func (uc *DreamsUseCase) Finalize(id uint) error {
	updates := map[string]any{
		"finalized": true,
	}
	return uc.repo.Update(id, updates)
}

func (uc *DreamsUseCase) Rate(id uint, rating int) error {
	updates := map[string]any{
		"rating": rating,
	}
	return uc.repo.Update(id, updates)
}
