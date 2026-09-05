package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration12 = &gormigrate.Migration{
	ID: "20260902_remove_visible",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct {
			Visible bool
		}
		return tx.Exec("ALTER TABLE dreams DROP COLUMN visible").Error
	},
	Rollback: func(tx *gorm.DB) error {
		type DreamVisible struct {
			Visible bool `gorm:"default:true"`
		}
		return tx.Table("dreams").AutoMigrate(&DreamVisible{})
	},
}
