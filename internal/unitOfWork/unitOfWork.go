package unitofwork

import (
	"context"

	"github.com/bbrighter/dreams-api/internal/categories"
	"github.com/bbrighter/dreams-api/internal/dreams"
	"gorm.io/gorm"
)

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Do(
	ctx context.Context,
	fn func(*dreams.DreamsRepo, *categories.CategoriesRepo) error,
) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoDreams := dreams.NewDreamsRepo(tx)
		repoCats := categories.NewCategoriesRepo(tx)
		return fn(repoDreams, repoCats)
	})
}
