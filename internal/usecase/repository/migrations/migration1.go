package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration1 = &gormigrate.Migration{
	ID: "20231218_Initial",
	Migrate: func(tx *gorm.DB) error {
		type Tag struct {
			ID      uint
			Title   string
			DreamID uint
		}
		type Dream struct {
			ID          uint
			Date        time.Time
			Description string
			Tags        []*Tag
		}

		return tx.AutoMigrate(
			Dream{},
			Tag{},
		)
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable("dreams"); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable("tags"); err != nil {
			return err
		}
		return nil
	},
}
