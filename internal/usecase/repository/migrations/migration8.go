package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration8 = &gormigrate.Migration{
	ID: "20250603_FixFinalizeDreams",
	Migrate: func(tx *gorm.DB) error {
		type Dream struct{}
		return tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&Dream{}).Update("finalized", gorm.Expr("NOT finalized")).Error

	},
	Rollback: func(tx *gorm.DB) error {
		return nil
	},
}
