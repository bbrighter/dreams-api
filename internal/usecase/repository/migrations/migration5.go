package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration5 = &gormigrate.Migration{
	ID: "20240218_HiddenDreams",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct {
			ID          uint
			Date        time.Time
			Description string
			Visible     *bool `gorm:"default:true"`
		}
		return tx.Migrator().AutoMigrate(&Dream{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropColumn("dreams", "visible")
	},
}
