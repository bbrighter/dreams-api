package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration4 = &gormigrate.Migration{
	ID: "20240207_RenameTags",
	Migrate: func(tx *gorm.DB) error {
		type Tag struct{}
		if err := tx.Migrator().RenameColumn(Tag{}, "title", "name"); err != nil {
			return err
		}
		if err := tx.Migrator().RenameTable("tags", "categories"); err != nil {
			return err
		}
		if err := tx.Migrator().RenameTable("tags_dreams", "categories_dreams"); err != nil {
			return err
		}
		type CategoriesDreams struct{}
		if err := tx.Migrator().RenameColumn(CategoriesDreams{}, "tag_id", "category_id"); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().RenameColumn("categories", "name", "title"); err != nil {
			return err
		}
		if err := tx.Migrator().RenameTable("categories", "tags"); err != nil {
			return err
		}
		if err := tx.Migrator().RenameTable("categories_dreams", "tags_dreams"); err != nil {
			return err
		}
		if err := tx.Migrator().RenameColumn("tags_dreams", "category_id", "tag_id"); err != nil {
			return err
		}
		return nil
	},
}
