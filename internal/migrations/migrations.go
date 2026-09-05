package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var migrations = []*gormigrate.Migration{
	migration1,
	migration2,
	migration3,
	migration4,
	migration5,
	migration6,
	migration7,
	migration8,
	migration9,
	migration10,
	migration11,
	migration12,
}

func migrationFactory(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}

func Migration(db *gorm.DB, logger *zap.Logger) {
	m := migrationFactory(db)
	if err := m.Migrate(); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}
	logger.Info("Migration successful")
}
