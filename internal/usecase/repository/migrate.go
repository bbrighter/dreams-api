package repository

import (
	mig "github.com/bbrighter/dreams-api/internal/usecase/repository/migrations"
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func migrationFactory(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, mig.Migrations)
}

func Migration(db *gorm.DB, logger *zap.Logger) {
	m := migrationFactory(db)
	if err := m.Migrate(); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}
	logger.Info("Migration successful")
}
