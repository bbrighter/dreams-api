package usecase

import (
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type DreamsUseCase struct {
	repo DreamsRepo
}

func New(r DreamsRepo) *DreamsUseCase {
	return &DreamsUseCase{
		repo: r,
	}
}

func (uc *DreamsUseCase) GetAll(showAll bool) entity.Dreams {
	return uc.repo.GetAll(showAll)
}

func (uc *DreamsUseCase) Get(id uint, showAll bool) (entity.Dream, error) {
	return uc.repo.GetById(id, showAll)
}

func (uc *DreamsUseCase) Create(date time.Time) (uint, error) {
	dream := entity.Dream{Date: date}
	return uc.repo.Create(dream)
}

func (uc *DreamsUseCase) Update(id uint, date time.Time, description string) error {
	dream := entity.Dream{ID: id, Date: date, Description: description}
	return uc.repo.Update(dream)
}

func (uc *DreamsUseCase) Delete(id uint) (entity.Categories, error) {
	dream := entity.Dream{ID: id}
	return uc.repo.Delete(dream)
}

func (uc *DreamsUseCase) ToggleVisibility(id uint) error {
	dream := entity.Dream{ID: id}
	return uc.repo.ToggleVisibility(dream)
}
