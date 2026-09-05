package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration6 = &gormigrate.Migration{
	ID: "20240728_VisibleNonPointer",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct {
			ID          uint
			Date        time.Time
			Description string
			Visible     bool `gorm:"default:true"`
		}
		return tx.AutoMigrate(&Dream{})
	},
	Rollback: func(tx *gorm.DB) error {
		type Dream struct {
			ID          uint
			Date        time.Time
			Description string
			Visible     *bool `gorm:"default:true"`
		}
		return tx.AutoMigrate(&Dream{})
	},
}
