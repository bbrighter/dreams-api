package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration2 = &gormigrate.Migration{
	ID: "20240105_Tags",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct{ ID uint }
		type Tag struct {
			ID     uint
			Title  string
			Dreams []Dream `gorm:"many2many:tags_dreams;"`
		}
		if err := tx.Migrator().DropTable("tags"); err != nil {
			return err
		}
		return tx.AutoMigrate(
			&Tag{},
		)
	},
	Rollback: func(tx *gorm.DB) error {
		type Tag struct {
			ID      uint
			Title   string
			DreamID uint
		}
		if err := tx.Migrator().DropTable("tags"); err != nil {
			return err
		}
		if err := tx.AutoMigrate(Tag{}); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable("tags_dreams"); err != nil {
			return err
		}
		return nil
	},
}
