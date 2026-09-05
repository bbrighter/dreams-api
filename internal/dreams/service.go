package dreams

import (
	"context"
	"time"

	apperrors "github.com/bbrighter/dreams-api/internal/appErrors"
	"github.com/bbrighter/dreams-api/internal/entities"
)

type DreamsService struct {
	d *DreamsRepo
}

func NewDreamsService(repo *DreamsRepo) *DreamsService {
	return &DreamsService{d: repo}
}

func (s *DreamsService) ListDreams(ctx context.Context) (entities.Dreams, error) {
	return s.d.ListDreams(ctx)
}

func (s *DreamsService) GetDreamById(ctx context.Context, id uint) (entities.Dream, error) {
	return s.d.GetDream(ctx, id)
}

func (s *DreamsService) CreateDream(ctx context.Context, date time.Time) (uint, error) {
	var dream = &entities.Dream{Date: date}
	return s.d.CreateDream(ctx, dream)
}

type UpdateDreamParams struct {
	Date        *time.Time
	Description *string
	Rating      *int
	Finalized   *bool
}

func (s *DreamsService) UpdateDream(ctx context.Context, id uint, params UpdateDreamParams) error {
	updates := make(map[string]any)
	if params.Date != nil {
		updates["date"] = *params.Date
	}
	if params.Description != nil {
		updates["description"] = *params.Description
	}
	if params.Rating != nil {
		updates["rating"] = *params.Rating
	}
	if params.Finalized != nil {
		updates["finalized"] = *params.Finalized
	}
	if len(updates) == 0 {
		return apperrors.ErrBadParam
	}

	if params.Finalized != nil && *params.Finalized {
		dream, err := s.GetDreamById(ctx, id)
		if err != nil {
			return err
		}
		if params.Rating == nil && dream.Rating == nil {
			return apperrors.ErrBadParam
		}
	}

	err := s.d.UpdateDream(ctx, id, updates)
	return err
}

func (s *DreamsService) DeleteDream(ctx context.Context, id uint) error {
	return s.d.DeleteDream(ctx, id)
}
