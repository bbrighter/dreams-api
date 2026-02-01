package usecase

import (
	"context"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"go.uber.org/zap"
)

type IDreams interface {
	List(context.Context) (entity.Dreams, error)
	Get(context.Context, uint) (entity.Dream, error)
	Create(context.Context, time.Time) (uint, error)
	Update(context.Context, uint, *time.Time, *string, *int) error
	Delete(context.Context, uint) error
	Finalize(context.Context, uint) error
}

type DreamsUseCase struct {
	BaseUseCase
	repo IDreamsRepo
}

func NewDreamUseCase(r IDreamsRepo, logger *zap.Logger) *DreamsUseCase {
	return &DreamsUseCase{
		repo: r, BaseUseCase: NewBaseUseCase(logger),
	}
}

func (uc *DreamsUseCase) List(ctx context.Context) (entity.Dreams, error) {
	return uc.repo.List(ctx, false)
}

func (uc *DreamsUseCase) Get(ctx context.Context, id uint) (entity.Dream, error) {
	dream, err := uc.repo.Get(ctx, id, false)
	return dream, err
}

func (uc *DreamsUseCase) Create(ctx context.Context, date time.Time) (uint, error) {
	dream := &entity.Dream{Date: date}
	id, err := uc.repo.Create(ctx, dream)
	return id, err
}

func (uc *DreamsUseCase) Update(ctx context.Context, id uint, date *time.Time, description *string, rating *int) error {
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
	err := uc.repo.Update(ctx, id, updates)
	return err
}

func (uc *DreamsUseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *DreamsUseCase) Finalize(ctx context.Context, id uint) error {
	dream, err := uc.repo.Get(ctx, id, true)
	if err != nil {
		return err
	}
	if dream.Rating == nil {
		return entity.ErrorBadParam
	}

	updates := map[string]any{
		"finalized": true,
	}
	return uc.repo.Update(ctx, id, updates)
}
