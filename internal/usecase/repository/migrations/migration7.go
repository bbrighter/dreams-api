package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration7 = &gormigrate.Migration{
	ID: "20250601_FinalizeDreams",
	Migrate: func(tx *gorm.DB) error {
		type Category struct{ ID uint }
		type Person struct{ ID uint }

		type Dream struct {
			ID          uint
			Date        time.Time
			Description string
			Visible     bool       `gorm:"default:true"`
			Categories  []Category `gorm:"many2many:categories_dreams;"`
			Persons     []Person   `gorm:"many2many:people_dreams;"`
			Finalized   bool
		}
		if err := tx.AutoMigrate(&Dream{}); err != nil {
			return err
		}
		var dreams []Dream
		err := tx.Preload("Categories").Preload("Persons").
			Where(`
					dreams.id NOT IN (
						SELECT dream_id FROM categories_dreams
					)
					OR dreams.id NOT IN (
						SELECT dream_id FROM people_dreams
					)
				`).Find(&dreams).Error
		if err != nil {
			return err
		}
		var ids []uint
		for _, d := range dreams {
			ids = append(ids, d.ID)
		}
		return tx.Model(&Dream{}).Where("id IN ?", ids).UpdateColumn("finalized", true).Error
	},
	Rollback: func(tx *gorm.DB) error {
		type Dream struct {
			ID          uint
			Date        time.Time
			Description string
			Visible     bool `gorm:"default:true"`
		}
		return tx.AutoMigrate(&Dream{})
	},
}
