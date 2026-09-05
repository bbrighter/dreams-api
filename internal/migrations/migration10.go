package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration10 = &gormigrate.Migration{
	ID: "20260619_Rating",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct {
			ID          uint
			Date        time.Time
			Description string
			Visible     bool `gorm:"default:true"`
			Finalized   bool
			Rating      *int
		}

		return tx.AutoMigrate(
			Dream{},
		)
	},
	Rollback: func(tx *gorm.DB) error {
		type Dream struct{ Rating *int }
		tx.Migrator().DropColumn(&Dream{}, "rating")
		return nil
	},
}
