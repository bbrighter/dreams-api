package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration3 = &gormigrate.Migration{
	ID: "20240131_Persons",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct{ ID uint }
		type Person struct {
			ID     uint
			Name   string
			Dreams []Dream `gorm:"many2many:people_dreams;"`
		}

		return tx.Migrator().AutoMigrate(&Person{})
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable("people_dreams"); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable("people"); err != nil {
			return err
		}
		return nil
	},
}
