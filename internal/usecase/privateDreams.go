package usecase

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type PrivateDreamsUseCase struct {
	repo IDreamsRepo
}

type PrivateDreams interface {
	ToggleVisibility(context.Context, uint) error
	List(context.Context) (entity.Dreams, error)
	Get(context.Context, uint) (entity.Dream, error)
}

func NewPrivateDreamUseCase(r IDreamsRepo) *PrivateDreamsUseCase {
	return &PrivateDreamsUseCase{repo: r}
}

func (uc PrivateDreamsUseCase) ToggleVisibility(ctx context.Context, id uint) error {
	return uc.repo.ToggleVisibility(ctx, id)
}

func (uc PrivateDreamsUseCase) List(ctx context.Context) (entity.Dreams, error) {
	return uc.repo.List(ctx, true)
}

func (uc PrivateDreamsUseCase) Get(ctx context.Context, id uint) (entity.Dream, error) {
	return uc.repo.Get(ctx, id, true)
}
