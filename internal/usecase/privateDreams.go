package usecase

import "github.com/bbrighter/dreams-api/internal/entity"

type PrivateDreamsUseCase struct {
	repo IDreamsRepo
}

func NewPrivateDreamUseCase(r IDreamsRepo) *PrivateDreamsUseCase {
	return &PrivateDreamsUseCase{repo: r}
}

func (uc PrivateDreamsUseCase) ToggleVisibility(id uint) error {
	var dream = entity.Dream{ID: id}
	return uc.repo.ToggleVisibility(dream)
}

func (uc PrivateDreamsUseCase) List() entity.Dreams {
	return uc.repo.List(true)
}

func (uc PrivateDreamsUseCase) Get(id uint) (entity.Dream, error) {
	return uc.repo.Get(id, true)
}
