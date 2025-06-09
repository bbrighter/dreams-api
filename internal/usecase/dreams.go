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

func (uc *DreamsUseCase) Update(id uint, date time.Time, description string) error {
	dream := entity.Dream{ID: id, Date: date, Description: description}
	return uc.repo.Update(dream)
}

func (uc *DreamsUseCase) Delete(id uint) (entity.Categories, entity.Persons, error) {
	dream := entity.Dream{ID: id}
	return uc.repo.Delete(dream)
}

func (uc *DreamsUseCase) Finalize(id uint) error {
	return uc.repo.Finalize(id)
}
